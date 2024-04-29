package middlewares

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
	
	"github.com/go-redis/redis/v8"
	"github.com/lucadboer/posgo/rate-limiter/internal/infra/storage"
)

type Limiter interface {
	Allow(key string) bool
}

type RateLimiter struct {
	Storage   storage.Storage
	Limit     int
	Window    time.Duration
	BlockTime time.Duration
	mu        sync.Mutex
}

func NewRateLimiter(storage storage.Storage, limit int, window, blockTime time.Duration) *RateLimiter {
	return &RateLimiter{
		Storage:   storage,
		Limit:     limit,
		Window:    window,
		BlockTime: blockTime,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	ctx := context.Background()
	
	blockedKey := fmt.Sprintf("%s:blocked", key)
	blockedVal, err := rl.Storage.Get(ctx, blockedKey)
	if err == nil && blockedVal == "blocked" {
		return false
	}
	
	val, err := rl.Storage.Get(ctx, key)
	
	if err == redis.Nil {
		if err := rl.Storage.Increment(ctx, key); err != nil {
			fmt.Println("Error incrementing key in storage:", err)
			return false
		}
		if err := rl.Storage.SetExpiration(ctx, key, rl.Window); err != nil {
			fmt.Println("Error setting expiration in storage:", err)
			return false
		}
		return true
	} else if err != nil {
		fmt.Println("Error accessing storage:", err)
		return false
	}
	
	count, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("Error converting value to int:", err)
		return false
	}
	
	if count > rl.Limit {
		if err := rl.Storage.SetBlock(ctx, blockedKey, rl.BlockTime); err != nil {
			fmt.Println("Error setting blocked status in storage:", err)
			return false
		}
		return false
	}
	
	if err = rl.Storage.Increment(ctx, key); err != nil {
		fmt.Println("Error incrementing key in storage:", err)
		return false
	}
	
	return true
}

func RateLimiterMiddleware(limiter Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				fmt.Println("Error extracting IP from RemoteAddr:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			
			if ip == "::1" {
				ip = "127.0.0.1"
			}
			
			key := ip
			if token := r.Header.Get("API_KEY"); token != "" {
				key = "token:" + token
			}
			
			if !limiter.Allow(key) {
				http.Error(w, "You have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}
