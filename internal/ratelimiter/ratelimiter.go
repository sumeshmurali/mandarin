package ratelimiter

import (
	"net/http"

	"github.com/sumeshmurali/mandarin/internal/config"
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

func NewRateLimiter(config *config.RatelimitConfig) Ratelimiter {
	switch config.RatelimitType {
	case "global":
		return &GlobalRatelimiter{
			limiter: rate.NewLimiter(rate.Limit(config.Ratelimit), config.Ratelimit),
		}
	}
	return nil
}

type GlobalRatelimiter struct {
	limiter *rate.Limiter
}

func (r *GlobalRatelimiter) Allow(_ *http.Request) bool {
	return r.limiter.Allow()
}

// type TokenEnforcedRatelimiter struct {
// 	limiter map[string]*rate.Limiter
// 	wlock   *sync.Mutex
// }

// func (r *TokenEnforcedRatelimiter) GetOrCreateLimiter(rq *http.Request) *rate.Limiter {
// 	r.wlock.Lock()
// 	defer r.wlock.Unlock()
// 	if _, ok := r.limiter[r.RemoteAddr]; !ok {
// 		r.limiter[r.RemoteAddr] = rate.NewLimiter(rate.Limit(1), 1)
// 	}
// 	return r.limiter[r.RemoteAddr]
// }
