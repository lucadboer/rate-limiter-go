package limiter

import (
	"context"
	"time"
)

type RateLimiter struct {
	store        Store
	ipLimit      int
	tokenLimit   int
	blockDuration time.Duration
}

func NewRateLimiter(store Store, ipLimit int, tokenLimit int, blockDuration int) *RateLimiter {
	return &RateLimiter{
		store:        store,
		ipLimit:      ipLimit,
		tokenLimit:   tokenLimit,
		blockDuration: time.Duration(blockDuration) * time.Second,
	}
}

func (rl *RateLimiter) isRateLimited(key string, limit int) bool {
	ctx := context.Background()
	// Incrementa o contador
	count, err := rl.store.Incr(ctx, key)
	if err != nil {
		return true // Em caso de erro, previna o tráfego excessivo
	}

	// Se for a primeira requisição, define o tempo de expiração
	if count == 1 {
		rl.store.Expire(ctx, key, rl.blockDuration)
	}

	return count > int64(limit)
}

func (rl *RateLimiter) IsIPRateLimited(ip string) bool {
	return rl.isRateLimited("ip:" + ip, rl.ipLimit)
}

func (rl *RateLimiter) IsTokenRateLimited(token string) bool {
	return rl.isRateLimited("token:" + token, rl.tokenLimit)
}
