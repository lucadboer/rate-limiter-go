package storage

import (
	"context"
	"time"
	
	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	Client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{Client: client}
}

func (rs *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	return rs.Client.Get(ctx, key).Result()
}

func (rs *RedisStorage) Increment(ctx context.Context, key string) error {
	_, err := rs.Client.Incr(ctx, key).Result()
	return err
}

func (rs *RedisStorage) SetExpiration(ctx context.Context, key string, duration time.Duration) error {
	_, err := rs.Client.Expire(ctx, key, duration).Result()
	return err
}

func (rs *RedisStorage) SetBlock(ctx context.Context, key string, duration time.Duration) error {
	_, err := rs.Client.Set(ctx, key, "blocked", duration).Result()
	return err
}
