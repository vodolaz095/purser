package mysql

import (
	"testing"

	"github.com/vodolaz095/purser/internal/repotest"
)

func TestRepo(t *testing.T) {
	mr := Repository{
		DatabaseConnectionString: "root:purser@tcp(127.0.0.1:3306)/purser?charset=utf8&parseTime=True&loc=Local",
	}
	repotest.ValidateRepo(t, "mariadb", &mr)
}
