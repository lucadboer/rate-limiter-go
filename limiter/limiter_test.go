package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func setupTestRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rdb.FlushAll(ctx)
	return rdb
}

func TestRateLimiterByIP(t *testing.T) {
	rdb := setupTestRedis()
	store := &RedisStore{Client: rdb}

	rl := NewRateLimiter(store, 2, 3, 10)

	assert.False(t, rl.IsIPRateLimited("192.168.1.1"))
	assert.False(t, rl.IsIPRateLimited("192.168.1.1"))
	assert.True(t, rl.IsIPRateLimited("192.168.1.1"))

	time.Sleep(11 * time.Second)
	assert.False(t, rl.IsIPRateLimited("192.168.1.1"))
}

func TestRateLimiterByToken(t *testing.T) {
	rdb := setupTestRedis()
	store := &RedisStore{Client: rdb}

	rl := NewRateLimiter(store, 2, 3, 10)

	assert.False(t, rl.IsTokenRateLimited("token123"))
	assert.False(t, rl.IsTokenRateLimited("token123"))
	assert.False(t, rl.IsTokenRateLimited("token123"))
	assert.True(t, rl.IsTokenRateLimited("token123"))

	time.Sleep(11 * time.Second)
	assert.False(t, rl.IsTokenRateLimited("token123"))
}

func TestRateLimiterPrioritizeTokenOverIP(t *testing.T) {
	rdb := setupTestRedis()
	store := &RedisStore{Client: rdb}

	rl := NewRateLimiter(store, 2, 3, 10)

	assert.False(t, rl.IsIPRateLimited("192.168.1.2"))
	assert.False(t, rl.IsIPRateLimited("192.168.1.2"))
	assert.True(t, rl.IsIPRateLimited("192.168.1.2"))

	assert.False(t, rl.IsTokenRateLimited("token456"))
	assert.False(t, rl.IsTokenRateLimited("token456"))
	assert.False(t, rl.IsTokenRateLimited("token456"))
	assert.True(t, rl.IsTokenRateLimited("token456"))
}
