package middleware

import (
	"context"
	"net/http"
	"time"
)

// TimeoutMiddleware is a middleware that sets a timeout for handling requests.
type TimeoutMiddleware struct {
	Next    http.Handler
	Timeout time.Duration // The maximum duration for request processing.
}

// NewTimeoutMiddleware creates a new TimeoutMiddleware instance.
func NewTimeoutMiddleware(next http.Handler, timeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{
		Next:    next,
		Timeout: timeout,
	}
}

// ServeHTTP is the middleware handler function that enforces the request timeout.
func (m *TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(r.Context(), m.Timeout)
	defer cancel()

	// Use the context with the timeout for handling the request.
	r = r.WithContext(ctx)

	// Call the next handler in the chain with the updated context.
	done := make(chan struct{})
	go func() {
		m.Next.ServeHTTP(w, r)
		close(done)
	}()

	// Wait for the request to finish or for the timeout to occur.
	select {
	case <-done:
		// The request completed within the timeout.
		return
	case <-ctx.Done():
		// The timeout has occurred.
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
	}
}
