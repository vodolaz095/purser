package watchdog

import (
	"context"
	"time"

	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/internal/service"
)

func Ready() (supported bool, err error) {
	return daemon.SdNotify(false, daemon.SdNotifyReady)
}

func StartWatchdog(ctx context.Context, ss *service.SecretService) {
	var err error
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil {
		return
	}
	if interval == 0 {
		log.Info().Msgf("Watchdog not enabled")
		return
	}
	ticker := time.NewTicker(interval / 2)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()

	for t := range ticker.C {
		pingCtx, cancel := context.WithDeadline(ctx, t.Add(interval/2))
		err = ss.Ping(pingCtx)
		if err != nil {
			_, err = daemon.SdNotify(false, daemon.SdNotifyWatchdog)
			if err != nil {
				log.Error().Err(err).Msgf("%s: while sending watchdog notification", err)
			} else {
				log.Debug().Msgf("SecretService is healthy!")
			}
		} else {
			log.Error().Err(err).Msgf("SecretService is broken! Ping error: %s", err)
		}
		cancel()
	}
}
