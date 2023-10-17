package memory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vodolaz095/purser/internal/repotest"
	"github.com/vodolaz095/purser/model"
)

var mr Repository

func TestRepo(t *testing.T) {
	repotest.ValidateRepo(t, "memory", &mr)

	mr.Broken = true
	err := mr.Ping(context.Background())
	if err != nil {
		if err.Error() != "service is broken" {
			t.Errorf("wrong error: %s", err)
		}
	} else {
		t.Errorf("error not thrown")
	}
}

func TestRepository_PutSecret(t *testing.T) {
	var err error
	err = mr.PutSecret(context.Background(), model.Secret{
		ID:   "a",
		Body: "aa",
		Meta: map[string]string{
			"aaa": "aaaa",
		},
		CreatedAt: time.Now().Add(time.Second - model.TTL),
		ExpireAt:  time.Now().Add(time.Second),
	})
	if err != nil {
		t.Errorf("ошибка задания секреа: %s", err)
	}
	found, err := mr.FindByID(context.Background(), "a")
	if err != nil {
		t.Errorf("ошибка поиска секрета: %s", err)
	}
	assert.Equal(t, "a", found.ID, "wrong id")
	assert.Equal(t, "aa", found.Body, "wrong body")
	assert.Equal(t, 1, len(found.Meta), "wrong meta length")
	assert.Equal(t, "aaaa", found.Meta["aaa"], "wrong meta")
	t.Logf("Secret will expire in %s", found.ExpireAt.Sub(time.Now()).String())
	t.Logf("Waiting for secret to expire...")
	time.Sleep(2 * time.Second)

	_, err = mr.FindByID(context.Background(), "a")
	if err != nil {
		if !errors.Is(err, model.ErrSecretNotFound) {
			t.Errorf("не та ошибка - %s", err)
		}
	} else {
		t.Errorf("ошибка не выдана при поиске не истекшего секрета")
	}
}

func TestRepository_Prune(t *testing.T) {
	err := mr.Prune(context.Background())
	if err != nil {
		t.Errorf("ошибка очистки: %s", err)
	}
	_, err = mr.FindByID(context.Background(), "a")
	if err != nil {
		if !errors.Is(err, model.ErrSecretNotFound) {
			t.Errorf("не та ошибка - %s", err)
		}
	} else {
		t.Errorf("ошибка не выдана при поиске не истекшего секрета")
	}
}
