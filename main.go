package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"

	"github.com/vodolaz095/purser/config"
	"github.com/vodolaz095/purser/internal/repository"
	"github.com/vodolaz095/purser/internal/repository/memory"
	"github.com/vodolaz095/purser/internal/service"
	purserGrpcServer "github.com/vodolaz095/purser/internal/transport/grpc"
	"github.com/vodolaz095/purser/internal/transport/grpc/proto"
	"github.com/vodolaz095/purser/internal/transport/watchdog"
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
		Str("hostname", config.Hostname).
		Str("environment", config.Environment).
		Logger().Level(zerolog.DebugLevel)
	log.Logger = sink

	log.Debug().Msgf("Application starting...")

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
	log.Debug().Msgf("Dialing jaeger on %s:%s", config.JaegerHost, config.JaegerHost)

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
	wg.Add(1)
	go func() {
		if config.ListenGRPC == "disabled" {
			return
		}
		listener, lErr := net.Listen("tcp", config.ListenGRPC)
		if err != nil {
			log.Error().Err(lErr).
				Msgf("error starting listener on %s: %s", config.ListenGRPC, lErr)
			return
		}
		grpcTransport := purserGrpcServer.PurserGrpcServer{
			Service: ss,
		}
		jwtMiddleware := purserGrpcServer.ValidateJWTInterceptor{HmacSecret: config.JwtSecret}
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				otelgrpc.UnaryServerInterceptor(),
				jwtMiddleware.ServerInterceptor,
			),
		))
		proto.RegisterPurserServer(grpcServer, &grpcTransport)
		log.Debug().Msgf("Preparing to start GRPC server on %s...", config.ListenGRPC)
		go func() {
			<-mainCtx.Done()
			log.Debug().Msg("Stopping GRPC server...")
			grpcServer.Stop()
			wg.Done()
		}()
		lErr = grpcServer.Serve(listener)
		if lErr != nil {
			log.Error().Err(err).Msgf("error starting grpc server on %s: %s", config.ListenGRPC, err)
		}
	}()

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
