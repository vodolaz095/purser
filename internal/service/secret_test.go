package service

import (
	"context"
	"os"
	"testing"

	"github.com/vodolaz095/purser/internal/repository"
	"github.com/vodolaz095/purser/internal/repository/memory"
	"github.com/vodolaz095/purser/internal/repository/mysql"
	"github.com/vodolaz095/purser/internal/repository/postgresql"
	"github.com/vodolaz095/purser/internal/repository/redis"
	"go.opentelemetry.io/otel"
)

func secretServiceTester(t *testing.T, repo repository.SecretRepo) {
	var err error
	ctx := context.TODO()
	err = repo.Init(ctx)
	if err != nil {
		t.Errorf("error initializing repo: %s", err)
		return
	}
	ss := SecretService{
		Tracer: otel.Tracer("unit_test_service"),
		Repo:   repo,
	}
	err = ss.Ping(ctx)
	if err != nil {
		t.Errorf("error pinging repo: %s", err)
		return
	}
	// TODO - other methods

	err = repo.Close(ctx)
	if err != nil {
		t.Errorf("error closing repo: %s", err)
	}
}

func TestSecretServiceMemory(t *testing.T) {
	repo := memory.Repo{}
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
