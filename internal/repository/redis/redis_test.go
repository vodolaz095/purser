package redis

import (
	"testing"

	"github.com/vodolaz095/purser/internal/repotest"
)

func TestRepo(t *testing.T) {
	rr := Repository{
		RedisConnectionString: "redis://127.0.0.1:6379",
	}
	repotest.ValidateRepo(t, "redis", &rr)
}
