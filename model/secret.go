package model

import (
	"errors"
	"time"
)

// https://gorm.io/docs/models.html

const TTL = 3 * time.Hour

var SecretNotFoundError = errors.New("secret not found")

type Secret struct {
	ID        string            `json:"id"`
	Body      string            `json:"body"`
	Meta      map[string]string `json:"fields"`
	CreatedAt time.Time         `json:"createdAt"`
	ExpireAt  time.Time         `json:"expireAt"`
}

func (s Secret) Expired() bool {
	return s.ExpireAt.Before(time.Now())
}
