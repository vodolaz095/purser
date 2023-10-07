package memory

import (
	"context"
	"sync"
	"time"

	"github.com/vodolaz095/purser/model"
	"github.com/vodolaz095/purser/pkg"
)

type Repo struct {
	sync.RWMutex
	data map[string]model.Secret
}

func (r *Repo) Init(ctx context.Context) error {
	r.data = make(map[string]model.Secret, 0)
	return nil
}

func (r *Repo) Ping(ctx context.Context) error {
	return nil
}

func (r *Repo) Close(ctx context.Context) error {
	r.Lock()
	r.data = nil
	r.Unlock()
	return nil
}

func (r *Repo) Create(ctx context.Context, body string, meta map[string]string) (model.Secret, error) {
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

func (r *Repo) FindByID(ctx context.Context, id string) (model.Secret, error) {
	r.RLock()
	defer r.RUnlock()
	secret, found := r.data[id]
	if found {
		if secret.Expired() {
			return model.Secret{}, model.SecretNotFoundError
		}
		return secret, nil
	}
	return model.Secret{}, model.SecretNotFoundError
}

func (r *Repo) DeleteByID(ctx context.Context, id string) error {
	r.Lock()
	defer r.Unlock()
	_, found := r.data[id]
	if found {
		delete(r.data, id)
		return nil
	}
	return model.SecretNotFoundError
}

func (r *Repo) Prune(ctx context.Context) error {
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
