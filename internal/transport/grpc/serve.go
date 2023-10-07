package grpc

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/internal/service"
	"github.com/vodolaz095/purser/internal/transport/grpc/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Options struct {
	HmacSecret string
	ListenOn   string
	Service    *service.SecretService
}

func Serve(ctx context.Context, opts Options) error {
	if opts.ListenOn == "disabled" {
		return nil
	}
	listener, err := net.Listen("tcp", opts.ListenOn)
	if err != nil {
		log.Error().Err(err).
			Msgf("error starting listener on %s: %s", opts.ListenOn, err)
		return err
	}
	grpcTransport := PurserGrpcServer{
		Service: *opts.Service,
	}
	jwtMiddleware := ValidateJWTInterceptor{HmacSecret: opts.HmacSecret}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			otelgrpc.UnaryServerInterceptor(),
			jwtMiddleware.ServerInterceptor,
		),
	))
	proto.RegisterPurserServer(grpcServer, &grpcTransport)
	log.Debug().Msgf("Preparing to start GRPC server on %s...", opts.ListenOn)
	go func() {
		<-ctx.Done()
		log.Debug().Msg("Stopping GRPC server...")
		grpcServer.Stop()
	}()
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error().Err(err).Msgf("error starting grpc server on %s: %s", opts.ListenOn, err)
		return err
	}
	return nil
}
