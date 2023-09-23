package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter is a middleware for rate limiting requests.
type RateLimiter struct {
	requestsPerSecond int
	bucket            chan struct{}
	mu                sync.Mutex
}

// NewRateLimiter creates a new RateLimiter with the specified requests per second limit.
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		bucket:            make(chan struct{}, requestsPerSecond),
	}

	// Fill the bucket initially.
	for i := 0; i < requestsPerSecond; i++ {
		rl.bucket <- struct{}{}
	}

	// Start a goroutine to replenish the bucket.
	go rl.startRefill()

	return rl
}

func (rl *RateLimiter) startRefill() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for i := len(rl.bucket); i < rl.requestsPerSecond; i++ {
			rl.bucket <- struct{}{}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	select {
	case <-rl.bucket:
		next(w, r)
	default:
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
	}
}
