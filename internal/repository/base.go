package repository

import "context"

// BaseRepo задаёт интерфейс, которому должны соответствовать все репозитории
type BaseRepo interface {
	// Ping проверяет соединение с базой данных
	Ping(ctx context.Context) error
	// Init настраивает соединение с базой данных
	Init(ctx context.Context) error
	// Close закрывает соединение с базой данных
	Close(ctx context.Context) error
}
