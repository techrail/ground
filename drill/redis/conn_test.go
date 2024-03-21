package redis

import (
	"github.com/techrail/ground/cache"
	"testing"
)

var client cache.Client

func TestRedisConnection(t *testing.T) {
	cache.CreateNewRedisClient()
}
