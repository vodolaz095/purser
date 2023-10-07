package http

import (
	"context"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/config"
	"github.com/vodolaz095/purser/internal/transport/http/middlewares"

	"github.com/vodolaz095/purser/internal/service"
)

type Options struct {
	HmacSecret string
	ListenOn   string
	Service    *service.SecretService
}

func Serve(ctx context.Context, opts Options) error {
	var err error
	if opts.ListenOn == "disabled" {
		return nil
	}
	listener, err := net.Listen("tcp", opts.ListenOn)
	if err != nil {
		log.Error().Err(err).
			Msgf("error starting listener on %s: %s", opts.ListenOn, err)
		return err
	}
	app := gin.New()
	if config.Environment != "production" {
		gin.SetMode(gin.DebugMode)
	}
	go func() {
		<-ctx.Done()
		err = listener.Close()
		if err != nil {
			log.Error().Err(err).
				Msgf("error closing http listener on %s : %s", opts.ListenOn, err)
		}
	}()
	middlewares.EmulatePHP(app)
	middlewares.UseTracing(app)
	middlewares.Secure(app)

	tr := Transport{
		Engine:  app,
		Service: *opts.Service,
	}

	tr.ExposeHealthChecks()
	tr.ExposeSecretAPI()

	log.Debug().Msgf("Preparing to start HTTP server on %s...", opts.ListenOn)
	err = app.RunListener(listener)
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return nil
		}
	}
	return err
}
