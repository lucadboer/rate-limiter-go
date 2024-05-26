package middleware

import (
	"net/http"

	"github.com/lucadboer/ratelimiter/limiter"
)

type RateLimiterMiddleware struct {
	Limiter *limiter.RateLimiter
}

func NewRateLimiterMiddleware(l *limiter.RateLimiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{Limiter: l}
}

func (rlm *RateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		token := r.Header.Get("API_KEY")
		
		if token != "" {
			if rlm.Limiter.IsTokenRateLimited(token) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		} else {
			if rlm.Limiter.IsIPRateLimited(ip) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
