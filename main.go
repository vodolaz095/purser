package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/config"
	"github.com/vodolaz095/purser/internal/repository"
	"github.com/vodolaz095/purser/internal/repository/memory"
	"github.com/vodolaz095/purser/internal/service"
	"github.com/vodolaz095/purser/internal/transport/watchdog"
	"go.opentelemetry.io/otel"

	"github.com/vodolaz095/purser/pkg"
)

func main() {
	wg := sync.WaitGroup{}
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error

	// logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}
	sink := zerolog.New(output).
		With().Timestamp().Caller().
		Logger().Level(zerolog.DebugLevel)
	log.Logger = sink

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
		Tracer: otel.Tracer("purser"),
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

	// start grpc server

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
	log.Debug().Msgf("Application is stopping")
}
