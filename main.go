package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
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

	// logging
	pkg.SetupLogger()

	// handle signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	go func() {
		s := <-sigc
		log.Info().Msgf("Signal %s is received", s.String())
		cancel()
	}()

	err = pkg.SetupJaeger(
		config.Hostname,
		config.Environment,
		config.JaegerHost,
		config.JaegerPort,
	)
	log.Debug().Msgf("Dialing jaeger on %s:%s", config.JaegerHost, config.JaegerPort)

	if err != nil {
		log.Fatal().Err(err).Msgf("error setting jaeger upd transfort for telemetry into %s:%s : %s",
			config.JaegerHost, config.JaegerPort, err)
	}

	// repository
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
		log.Fatal().Msgf("unknown database driver: %s", config.Driver)
	}
	err = repo.Init(mainCtx)
	if err != nil {
		log.Fatal().Err(err).Msgf("error pinging repo: %s", err)
	}
	log.Debug().Msgf("Repo initialized!")

	err = repo.Ping(mainCtx)
	if err != nil {
		log.Fatal().Err(err).Msgf("error pinging repo: %s", err)
	}
	log.Debug().Msgf("Repo online!")

	// service
	ss := service.SecretService{
		Tracer: otel.Tracer("purser_service_tracer"),
		Repo:   repo,
	}
	log.Debug().Msgf("Service initialized!")

	err = ss.Ping(mainCtx)
	if err != nil {
		log.Fatal().Err(err).Msgf("error pinging service: %s", err)
	}
	ss.Ready = true

	/*
	 * Transports
	 */

	// start systemd watchdog
	supported, err := watchdog.Ready()
	if err != nil {
		log.Fatal().Err(err).Msgf("error checking watchdog system: %s", err)
	}
	if supported {
		go watchdog.StartWatchdog(mainCtx, &ss)
	} else {
		log.Warn().Msgf("Systemd watchdog is disabled, application can be unstable")
	}
	// start http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		lErr := httpTransport.Serve(mainCtx, httpTransport.Options{
			HmacSecret: config.JwtSecret,
			ListenOn:   config.ListenHTTP,
			Service:    &ss,
		})
		if lErr != nil {
			log.Fatal().Err(lErr).Msgf("error starting http server on %s : %s",
				config.ListenHTTP, lErr)
		}
	}()
	// start grpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		lErr := grpcTransport.Serve(mainCtx, grpcTransport.Options{
			HmacSecret: config.JwtSecret,
			ListenOn:   config.ListenGRPC,
			Service:    &ss,
		})
		if lErr != nil {
			log.Fatal().Err(lErr).Msgf("error starting grpc server on %s : %s",
				config.ListenGRPC, lErr)
		}
	}()

	// start background routine to prune old secrets

	/*
	 * Shutdown properly
	 */

	// wait for main context to cancel
	wg.Add(1)
	go func() {
		<-mainCtx.Done()
		log.Info().Msgf("Closing main context...")
		errCloseRepo := repo.Close(context.Background())
		if errCloseRepo != nil {
			log.Error().Err(errCloseRepo).Msgf("Error closing repo: %s", errCloseRepo)
		}
		wg.Done()
	}()

	wg.Wait()
	log.Debug().Msgf("Application is stopped")
}
