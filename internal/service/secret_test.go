package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vodolaz095/purser/pkg/misc"
	"go.opentelemetry.io/otel"

	"github.com/vodolaz095/purser/internal/repository"
	"github.com/vodolaz095/purser/internal/repository/memory"
	"github.com/vodolaz095/purser/internal/repository/mysql"
	"github.com/vodolaz095/purser/internal/repository/postgresql"
	"github.com/vodolaz095/purser/internal/repository/redis"
	"github.com/vodolaz095/purser/model"
)

func secretServiceTester(t *testing.T, repo repository.SecretRepo) {
	var err error
	ctx := context.TODO()
	err = repo.Init(ctx)
	if err != nil {
		t.Errorf("error initializing repo: %s", err)
		return
	}
	defer func() {
		err = repo.Close(ctx)
		if err != nil {
			t.Errorf("error closing repo: %s", err)
		}
	}()
	ss := SecretService{
		Tracer: otel.Tracer("unit_test_service"),
		Repo:   repo,
	}
	err = ss.Ping(ctx)
	if err != nil {
		t.Errorf("error pinging repo: %s", err)
		return
	}
	secret, err := ss.Create(ctx, "test secret", map[string]string{
		"a": "b",
	})
	if err != nil {
		t.Errorf("error creating secret: %s", err)
		return
	}
	found, err := ss.FindByID(ctx, secret.ID)
	if err != nil {
		t.Errorf("error creating secret: %s", err)
		return
	}
	assert.Equal(t, secret.ID, found.ID, "id differs")
	assert.Equal(t, secret.Body, found.Body, "body differs")
	assert.Equal(t, secret.Meta, found.Meta, "meta differs")

	empty, err := ss.FindByID(ctx, misc.UUID())
	if err != nil {
		if err != model.ErrSecretNotFound {
			t.Errorf("wrong error: %s", err)
			return
		}
	}
	assert.Equal(t, empty.ID, "", "id differs")
	assert.Equal(t, empty.Body, "", "body differs")
	assert.Equal(t, len(empty.Meta), 0, "meta differs")

	err = ss.DeleteByID(ctx, secret.ID)
	if err != nil {
		t.Errorf("error deleting: %s", err)
		return
	}

	empty, err = ss.FindByID(ctx, secret.ID)
	if err != nil {
		if err != model.ErrSecretNotFound {
			t.Errorf("wrong error: %s", err)
			return
		}
	} else {
		t.Error("no error thrown for record not found")
	}
	assert.Equal(t, empty.ID, "", "id differs")
	assert.Equal(t, empty.Body, "", "body differs")
	assert.Equal(t, len(empty.Meta), 0, "meta differs")

	// проверяем хитрую бизнес логику
	programmersSecret, err := ss.Create(ctx, "Мне нравится язык программирования Golang", map[string]string{
		"a": "b",
	})
	if err != nil {
		t.Errorf("error creating secret: %s", err)
		return
	}
	found, err = ss.FindByID(ctx, programmersSecret.ID)
	if err != nil {
		t.Errorf("error creating secret: %s", err)
		return
	}
	assert.Equal(t, programmersSecret.ID, found.ID, "id differs")
	assert.Equal(t, programmersSecret.Body, found.Body, "body differs")
	assert.Equal(t, programmersSecret.Meta, found.Meta, "meta differs")
	assert.Equal(t, programmersSecret.Meta["a"], found.Meta["a"], "meta differs")
	assert.Equal(t, programmersSecret.Meta["programming"], found.Meta["programming"], "meta differs")
	assert.Equal(t, "yes", programmersSecret.Meta["programming"], "meta differs")
	assert.Equal(t, "yes", found.Meta["programming"], "meta differs")
}

func TestSecretServiceMemory(t *testing.T) {
	repo := memory.Repository{}
	secretServiceTester(t, &repo)
}

func TestSecretService_Prune(t *testing.T) {
	ctx := context.Background()
	repo := memory.Repository{}

	err := repo.PutSecret(ctx, model.Secret{
		ID:        "a",
		Body:      "aa",
		Meta:      map[string]string{"aaa": "aaaa"},
		CreatedAt: time.Now().Add(-model.TTL),
		ExpireAt:  time.Now(),
	})
	if err != nil {
		t.Errorf("error putting secret directly: %s", err)
		return
	}
	s := SecretService{
		Tracer: otel.Tracer("unit_test_service4prune"),
		Repo:   &repo,
	}
	err = s.Prune(ctx)
	if err != nil {
		t.Errorf("error pruning: %s", err)
	}
	_, err = repo.FindByID(ctx, "a")
	if err != nil {
		if !errors.Is(err, model.ErrSecretNotFound) {
			t.Errorf("не та ошибка - %s", err)
		}
	} else {
		t.Errorf("ошибка не выдана при поиске не истекшего секрета")
	}
}

func TestSecretServiceMysql(t *testing.T) {
	mysqlConnectionString := os.Getenv("MARIADB_DB_URL")
	if mysqlConnectionString != "" {
		repo := mysql.Repository{
			DatabaseConnectionString: mysqlConnectionString,
		}
		secretServiceTester(t, &repo)
	} else {
		t.Skipf("Переменная окружения MARIADB_DB_URL не задана - пропускаем тест")
	}
}

func TestSecretServiceRedis(t *testing.T) {
	redisConnectionString := os.Getenv("REDIS_DB_URL")
	if redisConnectionString != "" {
		repo := redis.Repository{
			RedisConnectionString: redisConnectionString,
		}
		secretServiceTester(t, &repo)
	} else {
		t.Skipf("Переменная окружения REDIS_DB_URL не задана - пропускаем тест")
	}
}

func TestSecretServicePostgresql(t *testing.T) {
	pgConString := os.Getenv("POSTGRES_DB_URL")
	if pgConString != "" {
		repo := postgresql.Repository{
			DatabaseConnectionString: pgConString,
		}
		secretServiceTester(t, &repo)
	} else {
		t.Skipf("Переменная окружения POSTGRES_DB_URL не задана - пропускаем тест")
	}
}
