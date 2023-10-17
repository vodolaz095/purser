package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vodolaz095/purser/model"
	"github.com/vodolaz095/purser/pkg/misc"
)

// Repository реализует интерфейс SecretRepo и
type Repository struct {
	sync.RWMutex
	data   map[string]model.Secret
	Broken bool
}

// Init настраивает соединение с базой данных
func (r *Repository) Init(_ context.Context) error {
	r.data = make(map[string]model.Secret, 0)
	return nil
}

// Ping проверяет соединение с базой данных
func (r *Repository) Ping(_ context.Context) error {
	if r.Broken {
		return fmt.Errorf("service is broken")
	}
	return nil
}

// Close закрывает соединение с базой данных
func (r *Repository) Close(_ context.Context) error {
	r.Lock()
	r.data = nil
	r.Unlock()
	return nil
}

// Create создаёт новый model.Secret
func (r *Repository) Create(_ context.Context, body string, meta map[string]string) (model.Secret, error) {
	r.Lock()
	defer r.Unlock()
	secret := model.Secret{
		ID:        misc.UUID(),
		Body:      body,
		Meta:      meta,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(model.TTL),
	}
	r.data[secret.ID] = secret
	return secret, nil
}

// FindByID ищет model.Secret по идентификатору
func (r *Repository) FindByID(_ context.Context, id string) (model.Secret, error) {
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
func (r *Repository) DeleteByID(_ context.Context, id string) error {
	r.Lock()
	defer r.Unlock()
	_, found := r.data[id]
	if found {
		delete(r.data, id)
		return nil
	}
	return model.ErrSecretNotFound
}

// PutSecret - не стандартный метод для тестирования, позволяющий задать секрет вручную во внутреннем хранилище.
// Используется для тестов.
func (r *Repository) PutSecret(_ context.Context, secret model.Secret) error {
	r.Lock()
	defer r.Unlock()
	if len(r.data) == 0 {
		r.data = make(map[string]model.Secret, 0)
	}
	r.data[secret.ID] = secret
	return nil
}

// Prune удаляет старые секреты
func (r *Repository) Prune(_ context.Context) error {
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
