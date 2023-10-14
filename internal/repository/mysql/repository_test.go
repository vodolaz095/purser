package mysql

import (
	"os"
	"testing"

	"github.com/vodolaz095/purser/internal/repotest"
)

func TestRepo(t *testing.T) {
	mysqlConnectionString := os.Getenv("MARIADB_DB_URL")
	if mysqlConnectionString != "" {
		repo := Repository{
			DatabaseConnectionString: mysqlConnectionString,
		}
		repotest.ValidateRepo(t, "mariadb", &repo)
	} else {
		t.Skipf("Переменная окружения MARIADB_DB_URL не задана - пропускаем тест")
	}
}
