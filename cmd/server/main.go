package main

import (
	"fmt"
	"github.com/lucadboer/posgo/rate-limiter/internal/http/middlewares"
	"github.com/lucadboer/posgo/rate-limiter/internal/infra/storage"
	"net/http"
	"time"
	
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
	configs "github.com/lucadboer/posgo/rate-limiter/config"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	
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
	
	rateLimiterWindow := 1 * time.Second
	rateLimiterBlockTime := 1 * time.Minute
	
	limiter := middlewares.NewRateLimiter(redisStorage, cfg.Limit, rateLimiterWindow, rateLimiterBlockTime)
	
	r.Use(middlewares.RateLimiterMiddleware(limiter))
	r.Use(middleware.Logger)
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	
	fmt.Println("Server is running...")
	
	http.ListenAndServe(":8080", r)
}
