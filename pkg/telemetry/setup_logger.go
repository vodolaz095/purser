package telemetry

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/config"
)

// SetupLogger внезапно, настраивает систему логгирования
func SetupLogger(hostname, environment, version string) {
	var output io.Writer
	switch config.LogOutputType(config.LogOutput) {
	case config.LogOutputConsole:
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}
		break
	case config.LogOutputStdOutJSON:
		output = os.Stdout
		break
	case config.LogOutputJournald:
		output = journald.NewJournalDWriter()
		break
	default:
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}
	}
	var level zerolog.Level
	switch config.LogLevel {
	case "trace":
		level = zerolog.TraceLevel
		break
	case "debug":
		level = zerolog.DebugLevel
		break
	case "info":
		level = zerolog.InfoLevel
		break
	case "warn":
		level = zerolog.WarnLevel
		break
	case "fatal":
		level = zerolog.FatalLevel
		break
	case "panic":
		level = zerolog.PanicLevel
		break
	default:
		level = zerolog.InfoLevel
	}

	sink := zerolog.New(output).
		With().Timestamp().Caller().
		Str("hostname", hostname).
		Str("environment", environment).
		Str("version", version).
		Logger().Level(level)
	log.Logger = sink

	log.Debug().Msgf("Логгирование настроено!")
}
