package storage

import (
	"context"
	"time"
)

type Storage interface {
	Get(ctx context.Context, key string) (string, error)
	Increment(ctx context.Context, key string) error
	SetExpiration(ctx context.Context, key string, duration time.Duration) error
	SetBlock(ctx context.Context, key string, duration time.Duration) error
}
