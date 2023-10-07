package service

import (
	"context"
	"testing"

	"github.com/vodolaz095/purser/internal/repository"
	"github.com/vodolaz095/purser/internal/repository/memory"
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
	t.Skipf("not implemented")
}

func TestSecretServiceRedis(t *testing.T) {
	t.Skipf("not implemented")
}

func TestSecretServicePostgresql(t *testing.T) {
	t.Skipf("not implemented")
}
