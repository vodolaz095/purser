package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vodolaz095/purser/model"
	"github.com/vodolaz095/purser/pkg"
)

// Repository реализует интерфейс SecretRepo с базой данных redis внутри
type Repository struct {
	RedisConnectionString string
	client                *redis.Client
}

// Ping проверяет соединение с базой данных
func (r *Repository) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Init настраивает соединение с базой данных
func (r *Repository) Init(ctx context.Context) error {
	opts, err := redis.ParseURL(r.RedisConnectionString)
	if err != nil {
		return err
	}
	r.client = redis.NewClient(opts)
	return r.Ping(ctx)
}

// Close закрывает соединение с базой данных
func (r *Repository) Close(ctx context.Context) error {
	return r.client.Close()
}

// Create создаёт новый model.Secret
func (r *Repository) Create(ctx context.Context, body string, meta map[string]string) (model.Secret, error) {
	id := pkg.UUID()
	meta["body"] = body
	pipe := r.client.Pipeline()
	for k := range meta {
		pipe.HSet(ctx, id, k, meta[k])
	}
	pipe.Expire(ctx, id, model.TTL)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return model.Secret{}, err
	}
	delete(meta, "body")
	return model.Secret{
		ID:        id,
		Body:      body,
		Meta:      meta,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(model.TTL),
	}, nil
}

// FindByID ищет model.Secret по идентификатору
func (r *Repository) FindByID(ctx context.Context, id string) (model.Secret, error) {
	var ret model.Secret
	raw, err := r.client.HGetAll(ctx, id).Result()
	if err != nil {
		return model.Secret{}, err
	}
	if len(raw) == 0 {
		return model.Secret{}, model.ErrSecretNotFound
	}
	ret.ID = id
	ret.Body = raw["body"]
	delete(raw, "body")
	ret.Meta = raw
	ttl, err := r.client.TTL(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return model.Secret{}, model.ErrSecretNotFound
		}
		return model.Secret{}, err
	}
	ret.ExpireAt = time.Now().Add(ttl)
	ret.CreatedAt = ret.ExpireAt.Add(-model.TTL) // lol
	return ret, nil
}

// DeleteByID удаляет секрет по идентификатору
func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	return r.client.Del(ctx, id).Err()
}

// Prune удаляет старые секреты
func (r *Repository) Prune(ctx context.Context) error {
	return nil
}
