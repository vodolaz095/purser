package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/internal/transport/prune"
	"go.opentelemetry.io/otel"

	"github.com/vodolaz095/purser/config"
	"github.com/vodolaz095/purser/internal/repository"
	"github.com/vodolaz095/purser/internal/repository/memory"
	"github.com/vodolaz095/purser/internal/repository/mysql"
	"github.com/vodolaz095/purser/internal/repository/postgresql"
	"github.com/vodolaz095/purser/internal/repository/redis"
	"github.com/vodolaz095/purser/internal/service"
	grpcTransport "github.com/vodolaz095/purser/internal/transport/grpc"
	httpTransport "github.com/vodolaz095/purser/internal/transport/http"
	"github.com/vodolaz095/purser/internal/transport/watchdog"
	"github.com/vodolaz095/purser/pkg"
)

func main() {
	wg := sync.WaitGroup{}
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error

	// настраиваем логгирование
	pkg.SetupLogger()

	// настраиваем приём сигналов от операционной системы сигналы
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	go func() {
		s := <-sigc
		log.Info().Msgf("Получен сигнал %s от операционной системы...", s.String())
		cancel()
	}()

	// настраиваем соединение с приёмником телеметрии
	log.Debug().Msgf("Соединяемся с сервисом телеметрии по %s:%s", config.JaegerHost, config.JaegerPort)
	err = pkg.SetupJaeger(
		config.Hostname,
		config.Environment,
		config.JaegerHost,
		config.JaegerPort,
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("Ошибка соединяемся с сервисом телеметрии по %s:%s : %s",
			config.JaegerHost, config.JaegerPort, err)
	}

	/*
	 * Настраиваем репозиторий для объектов типа model.Secret
	 */
	var repo repository.SecretRepo
	switch config.Driver {
	case "memory":
		repo = &memory.Repo{}
		break
	case "redis":
		repo = &redis.Repository{RedisConnectionString: config.DatabaseConnectionString}
		break
	case "mariadb", "mysql":
		repo = &mysql.Repository{DatabaseConnectionString: config.DatabaseConnectionString}
		break
	case "postgres", "pgx":
		repo = &postgresql.Repository{DatabaseConnectionString: config.DatabaseConnectionString}
		break
	default:
		log.Fatal().Msgf("неизвестный драйвер базы данных для репозитория: %s", config.Driver)
	}
	err = repo.Init(mainCtx)
	if err != nil {
		log.Fatal().Err(err).Msgf("ошибка инициализации репозитория: %s", err)
	}
	log.Debug().Msgf("Репозиторий инициализирован!")

	err = repo.Ping(mainCtx)
	if err != nil {
		log.Fatal().Err(err).Msgf("ошибка проверки репозитория: %s", err)
	}
	log.Debug().Msgf("Репозиторий готов к работе!")

	/*
	 * Настраиваем сервисы
	 */
	cs := service.CounterService{}
	cs.Init()

	ss := service.SecretService{
		Tracer: otel.Tracer("purser_service_tracer"),
		Repo:   repo,
	}
	log.Debug().Msgf("Сервис секретов инициализирован!")

	err = ss.Ping(mainCtx)
	if err != nil {
		log.Fatal().Err(err).Msgf("ошибка проверки сервиса секретов: %s", err)
	}
	/*
	 * Настраиваем транспорты
	 */

	// запускаем systemd watchdog который будет проверять корректность работы сервиса под управлением systemd
	supported, err := watchdog.Ready()
	if err != nil {
		log.Fatal().Err(err).Msgf("ошибка проверки Systemd Watchdog : %s", err)
	}
	if supported {
		go watchdog.StartWatchdog(mainCtx, &ss)
	} else {
		log.Warn().Msgf("Systemd Watchdog не активирован, работа приложения может быть нестабильной")
	}
	// Запускаем HTTP сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		lErr := httpTransport.Serve(mainCtx, httpTransport.Options{
			HmacSecret:     config.JwtSecret,
			ListenOn:       config.ListenHTTP,
			SecretService:  &ss,
			CounterService: &cs,
		})
		if lErr != nil {
			log.Fatal().Err(lErr).Msgf("Ошибка запуска HTTP сервера на %s : %s",
				config.ListenHTTP, lErr)
		}
	}()
	// Запускаем gRPC сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		lErr := grpcTransport.Serve(mainCtx, grpcTransport.Options{
			HmacSecret:     config.JwtSecret,
			ListenOn:       config.ListenGRPC,
			SecretService:  &ss,
			CounterService: &cs,
		})
		if lErr != nil {
			log.Fatal().Err(lErr).Msgf("Ошибка запуска gRPC сервера на %s : %s",
				config.ListenGRPC, lErr)
		}
	}()

	// запускаем фоновой процесс очистки старых документов
	av := prune.Autovacuum{Service: ss}
	wg.Add(1)
	go func() {
		av.StartPruningExpiredSecrets(mainCtx, config.PruneOldSecretsInterval)
		wg.Done()
	}()

	// Ждём, как контекст завершится, чтобы правильно закрыть репозиторий
	wg.Add(1)
	go func() {
		<-mainCtx.Done()
		log.Debug().Msgf("Завершение главного контекста приложения")
		errCloseRepo := repo.Close(context.Background())
		if errCloseRepo != nil {
			log.Error().Err(errCloseRepo).Msgf("Ошибка закрытия репозитория: %s", errCloseRepo)
		}
		log.Debug().Msgf("Репозиторий закрыт")
		wg.Done()
	}()
	wg.Wait()
	log.Info().Msgf("Сервис остановлен штатно")
}
