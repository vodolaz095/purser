package prune

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/internal/service"
)

// Timeout задаёт допустимую длительность удаления старых записей
const Timeout = 5 * time.Second

// Autovacuum реализует транспорт, который удаляет старые записи по таймеру
type Autovacuum struct {
	Service service.SecretService
}

// StartPruningExpiredSecrets запускает удаление старых записей по таймеру
func (av *Autovacuum) StartPruningExpiredSecrets(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("Останавливаем таймер очистки")
			ticker.Stop()
			return
		case <-ticker.C:
			ctx2, cancel := context.WithTimeout(ctx, Timeout)
			err := av.Service.Ping(ctx2)
			if err != nil {
				log.Error().Err(err).
					Msgf("Ошибка проверки статуса сервиса: %s", err)
				cancel()
				continue // вдруг обойдётся
			}
			err = av.Service.Prune(ctx2)
			if err != nil {
				log.Error().Err(err).
					Msgf("Ошибка очистки: %s", err)
			} else {
				log.Debug().Msgf("Старые секреты удалены!")
			}
			cancel()
			break
		}
	}
}
