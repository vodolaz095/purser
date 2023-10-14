package postgresql

import (
	"os"
	"testing"

	"github.com/vodolaz095/purser/internal/repotest"
)

func TestRepo(t *testing.T) {
	pgConString := os.Getenv("POSTGRES_DB_URL")
	if pgConString != "" {
		mr := Repository{
			DatabaseConnectionString: pgConString,
		}
		repotest.ValidateRepo(t, "pgx", &mr)
	} else {
		t.Skipf("Переменная окружения POSTGRES_DB_URL не задана - пропускаем тест")
	}
}
