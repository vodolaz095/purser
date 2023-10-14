package memory

import (
	"context"
	"testing"

	"github.com/vodolaz095/purser/internal/repotest"
)

func TestRepo(t *testing.T) {
	var mr Repository
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
