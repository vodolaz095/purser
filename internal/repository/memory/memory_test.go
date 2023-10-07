package memory

import (
	"testing"

	"github.com/vodolaz095/purser/internal/repotest"
)

func TestRepo(t *testing.T) {
	var mr Repo
	repotest.ValidateRepo(t, "memory", &mr)
}
