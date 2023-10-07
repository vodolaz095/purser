package repository

import (
	"context"

	"github.com/vodolaz095/purser/model"
)

type SecretRepo interface {
	BaseRepo
	Create(ctx context.Context, body string, meta map[string]string) (model.Secret, error)
	FindByID(ctx context.Context, id string) (model.Secret, error)
	DeleteByID(ctx context.Context, id string) error
	Prune(context.Context) error
}
