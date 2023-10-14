package repository

import (
	"context"

	"github.com/vodolaz095/purser/model"
)

// SecretRepo задаёт интерфейс, которому должен соответствовать репозиторий для работы с model.Secret
type SecretRepo interface {
	BaseRepo
	// Create создаёт новый секрет
	Create(ctx context.Context, body string, meta map[string]string) (model.Secret, error)
	// FindByID ищет секрет по идентификатору
	FindByID(ctx context.Context, id string) (model.Secret, error)
	// DeleteByID удаляет секрет по идентификатору
	DeleteByID(ctx context.Context, id string) error
	// Prune удаляет все устаревшие секреты
	Prune(context.Context) error
}
