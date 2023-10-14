package http

import (
	"context"
	"net"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/config"
	"github.com/vodolaz095/purser/internal/transport/http/middlewares"

	"github.com/vodolaz095/purser/internal/service"
)

type Options struct {
	HmacSecret     string
	ListenOn       string
	Hostname       string
	SecretService  *service.SecretService
	CounterService *service.CounterService
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
	} else {
		app.TrustedPlatform = gin.PlatformCloudflare
	}
	go func() {
		<-ctx.Done()
		log.Debug().Msg("Останавливаем HTTP сервер...")
		err = listener.Close()
		if err != nil {
			log.Error().Err(err).
				Msgf("error closing http listener on %s : %s", opts.ListenOn, err)
		}
		log.Debug().Msg("HTTP сервер остановлен")
	}()
	app.Use(
		middlewares.EmulatePHP(),
		middlewares.UseTracing(),
		middlewares.Secure(),
		middlewares.AddPermissionPolicyHeader(),
	)
	tr := Transport{
		Engine:         app,
		Hostname:       opts.Hostname,
		SecretService:  opts.SecretService,
		CounterService: opts.CounterService,
	}

	tr.ExposeHealthChecks()
	tr.ExposeSecretAPI()
	tr.ExposeMetrics()

	if !config.IsProduction() {
		log.Warn().Msgf("Система удалённой отладки pprof доступна по /debug/pprof!")
		pprof.Register(app)
	}

	log.Debug().Msgf("Preparing to start HTTP server on %s...", opts.ListenOn)
	err = app.RunListener(listener)
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return nil
		}
	}
	return err
}
