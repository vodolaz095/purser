package repository

import "context"

type BaseRepo interface {
	Ping(ctx context.Context) error
	Init(ctx context.Context) error
	Close(ctx context.Context) error
}
