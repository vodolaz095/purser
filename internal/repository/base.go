package repository

import "context"

// BaseRepo задаёт интерфейс, которому должны соответствовать все репозитории
type BaseRepo interface {
	Ping(ctx context.Context) error
	Init(ctx context.Context) error
	Close(ctx context.Context) error
}
