package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/lab42/httplib/internal"
	log "github.com/rs/zerolog/log"
)

// ZeroLogMiddleware is a middleware that logs HTTP request information using ZeroLog.
type ZeroLogMiddleware struct {
	Next          http.Handler
	SensitiveKeys []string // Sensitive field keys to be obfuscated
}

// NewZeroLogMiddleware creates a new ZeroLogMiddleware instance.
func NewZeroLogMiddleware(next http.Handler, sensitiveKeys []string) *ZeroLogMiddleware {
	return &ZeroLogMiddleware{Next: next, SensitiveKeys: sensitiveKeys}
}

// ServeHTTP is the middleware handler function that logs HTTP request information.
func (m *ZeroLogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Record the start time to calculate the request duration
	now := time.Now()

	// Create a ZeroLog event for logging
	event := log.Info()

	// Log query parameters
	for k, v := range r.URL.Query() {
		event.Str("query_"+k, strings.Join(v, ","))
	}

	// Parse form data and log it
	r.ParseForm()
	for k, v := range r.Form {
		// Obfuscate sensitive data based on configured sensitive keys
		obfuscatedValues := internal.ObfuscateSensitiveData(k, v, m.SensitiveKeys)
		event.Str("form_"+k, strings.Join(obfuscatedValues, ","))
	}

	// Log request details
	event.Str("method", r.Method)
	event.Str("protocol", r.Proto)
	event.Str("remote_addr", r.RemoteAddr)
	event.Str("request_uri", r.RequestURI)
	event.Str("user_agent", r.UserAgent())
	event.Str("referer", r.Referer())

	// Call the next handler in the chain
	m.Next.ServeHTTP(w, r)

	// Calculate and log the request duration in milliseconds
	duration := time.Since(now).Milliseconds()
	event.Int64("duration_ms", duration)
}
