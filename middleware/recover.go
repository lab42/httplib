package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// RecoveryMiddleware is a middleware that recovers from panics and logs errors.
type RecoveryMiddleware struct {
	Next http.Handler
}

// NewRecoveryMiddleware creates a new RecoveryMiddleware instance.
func NewRecoveryMiddleware(next http.Handler) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		Next: next,
	}
}

// ServeHTTP is the middleware handler function that recovers from panics and logs errors.
func (m *RecoveryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			// Recover from the panic.
			fmt.Println("Panic recovered:", r)

			// Log the stack trace for debugging purposes.
			debug.PrintStack()

			// Respond with an internal server error.
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	// Call the next handler in the chain.
	m.Next.ServeHTTP(w, r)
}
