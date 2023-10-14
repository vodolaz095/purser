package redis

import (
	"os"
	"testing"

	"github.com/vodolaz095/purser/internal/repotest"
)

func TestRepo(t *testing.T) {
	redisConnectionString := os.Getenv("REDIS_DB_URL")
	if redisConnectionString != "" {
		rr := Repository{
			RedisConnectionString: redisConnectionString,
		}
		repotest.ValidateRepo(t, "redis", &rr)
	} else {
		t.Skipf("Переменная окружения REDIS_DB_URL не задана - пропускаем тест")
	}

}
