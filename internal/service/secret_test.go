package service

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"

	"github.com/vodolaz095/purser/internal/repository"
	"github.com/vodolaz095/purser/internal/repository/memory"
	"github.com/vodolaz095/purser/internal/repository/mysql"
	"github.com/vodolaz095/purser/internal/repository/postgresql"
	"github.com/vodolaz095/purser/internal/repository/redis"
	"github.com/vodolaz095/purser/model"
	"github.com/vodolaz095/purser/pkg"
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

	empty, err := ss.FindByID(ctx, pkg.UUID())
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
}

func TestSecretServiceMemory(t *testing.T) {
	repo := memory.Repository{}
	secretServiceTester(t, &repo)
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
