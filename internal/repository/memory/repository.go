package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vodolaz095/purser/model"
	"github.com/vodolaz095/purser/pkg"
)

// Repository реализует интерфейс SecretRepo и
type Repository struct {
	sync.RWMutex
	data   map[string]model.Secret
	Broken bool
}

// Init настраивает соединение с базой данных
func (r *Repository) Init(ctx context.Context) error {
	r.data = make(map[string]model.Secret, 0)
	return nil
}

// Ping проверяет соединение с базой данных
func (r *Repository) Ping(ctx context.Context) error {
	if r.Broken {
		return fmt.Errorf("service is broken")
	}
	return nil
}

// Close закрывает соединение с базой данных
func (r *Repository) Close(ctx context.Context) error {
	r.Lock()
	r.data = nil
	r.Unlock()
	return nil
}

// Create создаёт новый model.Secret
func (r *Repository) Create(ctx context.Context, body string, meta map[string]string) (model.Secret, error) {
	r.Lock()
	defer r.Unlock()
	secret := model.Secret{
		ID:        pkg.UUID(),
		Body:      body,
		Meta:      meta,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(model.TTL),
	}
	r.data[secret.ID] = secret
	return secret, nil
}

// FindByID ищет model.Secret по идентификатору
func (r *Repository) FindByID(ctx context.Context, id string) (model.Secret, error) {
	r.RLock()
	defer r.RUnlock()
	secret, found := r.data[id]
	if found {
		if secret.Expired() {
			return model.Secret{}, model.ErrSecretNotFound
		}
		return secret, nil
	}
	return model.Secret{}, model.ErrSecretNotFound
}

// DeleteByID удаляет секрет по идентификатору
func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	r.Lock()
	defer r.Unlock()
	_, found := r.data[id]
	if found {
		delete(r.data, id)
		return nil
	}
	return model.ErrSecretNotFound
}

// Prune удаляет старые секреты
func (r *Repository) Prune(ctx context.Context) error {
	r.Lock()
	defer r.Unlock()
	keysToDelete := make([]string, 0, len(r.data)) // just in case all records expired
	for k := range r.data {
		if r.data[k].Expired() {
			keysToDelete = append(keysToDelete, k)
		}
	}
	for k := range keysToDelete {
		delete(r.data, keysToDelete[k])
	}
	return nil
}
