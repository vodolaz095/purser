package repotest

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vodolaz095/purser/internal/repository"
	"github.com/vodolaz095/purser/model"
	"github.com/vodolaz095/purser/pkg/misc"
)

// ValidateRepo используется в юнит тестах, чтобы базово проверить репозиторий
func ValidateRepo(t *testing.T, name string, repo repository.SecretRepo) {
	ctx := context.TODO()
	var err error
	defer func() {
		err = repo.Close(ctx)
		if err != nil {
			t.Errorf("error closing repo : %v", err)
			return
		}
		t.Log("repo is closed properly")
	}()

	err = repo.Init(ctx)
	if err != nil {
		t.Errorf("error initializing repo : %v", err)
		return
	}
	t.Logf("Repo initialized")
	err = repo.Ping(ctx)
	if err != nil {
		t.Errorf("error pinging repo : %v", err)
		return
	}
	t.Logf("Repo pinged")
	secret, err := repo.Create(ctx,
		fmt.Sprintf("test body for repo %s", name),
		map[string]string{
			"repo": name,
		},
	)
	if err != nil {
		t.Errorf("error creating secret : %v", err)
		return
	}
	t.Logf("Secret ID is %s", secret.ID)
	t.Logf("Secret will expire in %s", secret.ExpireAt.Sub(time.Now()).String())
	t.Logf("Repo %s allows to create object", name)

	secretExtracted, err := repo.FindByID(ctx, secret.ID)
	if err != nil {
		t.Errorf("error finding secret : %v", err)
		return
	}
	assert.Equal(t, secret.ID, secretExtracted.ID, "id differs")
	assert.Equal(t, secret.Body, secretExtracted.Body, "body differs")
	assert.Equal(t, len(secret.Meta), len(secretExtracted.Meta), "meta size differs")
	for k := range secret.Meta {
		assert.Equalf(t, secret.Meta[k], secretExtracted.Meta[k], "meta %s differs", k)
	}
	for k := range secretExtracted.Meta {
		assert.Equalf(t, secretExtracted.Meta[k], secret.Meta[k], "meta %s differs", k)
	}
	t.Logf("Repo %s returns proper secret by known key", name)

	unknownID := misc.UUID()

	secretNotFound, err := repo.FindByID(ctx, unknownID)
	if err != nil {
		if errors.Is(err, model.ErrSecretNotFound) {
			t.Logf("expected error returned for secret not found")
		} else {
			t.Errorf("error finding secret : %v", err)
			return
		}
	} else {
		t.Error("error not thrown for secret not found")
		return
	}
	assert.Equal(t, "", secretNotFound.ID, "not found secret's id is not null")
	assert.Equal(t, "", secretNotFound.Body, "not found secret's body is not null")
	assert.Empty(t, secretNotFound.Meta, "not found secret's meta is not null")
	t.Logf("Repo %s returns proper error for secret not found", name)

	err = repo.DeleteByID(ctx, secret.ID)
	if err != nil {
		t.Errorf("error deleting existent secret : %v", err)
		return
	}
	t.Logf("Repo allows to delete secret")
	secretThatShouldBeNotFound, err := repo.FindByID(ctx, secret.ID)
	if err != nil {
		if errors.Is(err, model.ErrSecretNotFound) {
			t.Logf("expected error returned for secret not found")
		} else {
			t.Errorf("error finding secret : %v", err)
			return
		}
	} else {
		t.Error("error not thrown for secret not found")
	}
	assert.Equal(t, "", secretThatShouldBeNotFound.ID, "not found secret's id is not null")
	assert.Equal(t, "", secretThatShouldBeNotFound.Body, "not found secret's body is not null")
	assert.Empty(t, secretThatShouldBeNotFound.Meta, "not found secret's meta is not null")
	t.Logf("Repo %s allows secret to be deleted", name)
}
