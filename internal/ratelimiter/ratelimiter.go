package ratelimiter

import (
	"net/http"

	"golang.org/x/time/rate"
)

type Ratelimiter interface {
	Allow(r *http.Request) bool
}

func RatelimitedHandlerMiddleWare(rl Ratelimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rl.Allow(r) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	}
}
func RatelimitedHandlerMiddleWareCurry(rl Ratelimiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		if rl == nil {
			// if no ratelimiter, return the next handler
			return next
		}
		return RatelimitedHandlerMiddleWare(rl, next)
	}
}

type GlobalRatelimiter struct {
	limiter *rate.Limiter
}

func (r *GlobalRatelimiter) Allow(_ *http.Request) bool {
	return r.limiter.Allow()
}

func NewRateLimiter(ratelimiterType string, ratelimit int) Ratelimiter {
	switch ratelimiterType {
	case "global":
		return &GlobalRatelimiter{
			limiter: rate.NewLimiter(rate.Limit(ratelimit), ratelimit),
		}
	}
	return nil
}
