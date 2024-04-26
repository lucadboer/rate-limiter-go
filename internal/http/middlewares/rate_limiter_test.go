package middlewares

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"testing"
	"time"
)

type MockStorage struct {
	Data map[string]string
}

func (m *MockStorage) Get(ctx context.Context, key string) (string, error) {
	val, exists := m.Data[key]
	if !exists {
		return "", redis.Nil
	}
	return val, nil
}

func (m *MockStorage) Increment(ctx context.Context, key string) error {
	val, exists := m.Data[key]
	if !exists {
		m.Data[key] = "1"
	} else {
		count, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		m.Data[key] = strconv.Itoa(count + 1)
	}
	return nil
}

func (m *MockStorage) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (m *MockStorage) SetBlock(ctx context.Context, key string, duration time.Duration) error {
	m.Data[key] = "blocked"
	return nil
}

func TestAllow(t *testing.T) {
	mockStorage := &MockStorage{Data: make(map[string]string)}
	limiter := NewRateLimiter(mockStorage, 10, 1*time.Second, 1*time.Minute)
	
	for i := 0; i < 11; i++ {
		if !limiter.Allow("testKey") {
			t.Errorf("Allow was false, want true for request number %d", i+1)
		}
	}
	
	if limiter.Allow("testKey") {
		t.Errorf("Allow was true, want false for request number 6")
	}
}
