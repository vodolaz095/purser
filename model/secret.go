package model

import (
	"errors"
	"time"
)

// TTL задаёт срок жизни секрета
const TTL = 3 * time.Hour

// SecretNotFoundError ошибка, возвращаемая, если секрет не найден в хранилище
var SecretNotFoundError = errors.New("secret not found")

// Secret - структура данных с которой работает приложение
type Secret struct {
	ID        string            `json:"id"`
	Body      string            `json:"body"`
	Meta      map[string]string `json:"fields"`
	CreatedAt time.Time         `json:"createdAt"`
	ExpireAt  time.Time         `json:"expireAt"`
}

// Expired проверяет, устарел ни секрет
func (s Secret) Expired() bool {
	return s.ExpireAt.Before(time.Now())
}
