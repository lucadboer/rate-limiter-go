package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPassword           string
	RateLimit            int
	RateLimiterWindow    time.Duration
	RateLimiterBlockTime time.Duration
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	limit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		limit = 10
	}

	window, err := time.ParseDuration(os.Getenv("RATE_LIMITER_WINDOW"))
	if err != nil {
		window = 1 * time.Second
	}

	blockTime, err := time.ParseDuration(os.Getenv("RATE_LIMITER_BLOCK_TIME"))
	if err != nil {
		blockTime = 1 * time.Minute
	}

	return &Config{
		DBPassword:           os.Getenv("DB_PASSWORD"),
		RateLimit:            limit,
		RateLimiterWindow:    window,
		RateLimiterBlockTime: blockTime,
	}, nil
}
