package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
	configs "github.com/lucadboer/posgo/rate-limiter/config"
	"github.com/lucadboer/posgo/rate-limiter/internal/http/middlewares"
	"github.com/lucadboer/posgo/rate-limiter/internal/infra/storage"
	"net/http"
)

func main() {
	cfg, err := configs.LoadConfig()

	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	redisStorage := storage.NewRedisStorage(redisClient)

	r := chi.NewRouter()

	limiter := middlewares.NewRateLimiter(redisStorage, cfg.RateLimit, cfg.RateLimiterWindow, cfg.RateLimiterBlockTime)

	r.Use(middlewares.RateLimiterMiddleware(limiter))
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	fmt.Println("Server is running...")

	http.ListenAndServe(":8080", r)
}
