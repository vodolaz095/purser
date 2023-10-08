package postgresql

import (
	"testing"

	"github.com/vodolaz095/purser/internal/repotest"
)

func TestRepo(t *testing.T) {
	mr := Repository{
		DatabaseConnectionString: "postgres://purser:purser@127.0.0.1:5432/purser",
	}
	repotest.ValidateRepo(t, "pgx", &mr)
}
