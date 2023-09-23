package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/lab42/httplib/internal"
	"go.uber.org/zap"
)

// ZapMiddleware is a middleware that logs HTTP request information using Zap.
type ZapMiddleware struct {
	Next          http.Handler
	Logger        *zap.Logger
	SensitiveKeys []string // Sensitive field keys to be obfuscated
}

// NewZapMiddleware creates a new ZapMiddleware instance.
func NewZapMiddleware(next http.Handler, logger *zap.Logger, sensitiveKeys []string) *ZapMiddleware {
	return &ZapMiddleware{Next: next, Logger: logger, SensitiveKeys: sensitiveKeys}
}

// ServeHTTP is the middleware handler function that logs HTTP request information.
func (m *ZapMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Record the start time to calculate the request duration
	now := time.Now()

	// Create a Zap logger for logging
	logger := m.Logger

	// Log query parameters
	for k, v := range r.URL.Query() {
		logger.Info("Query parameter",
			zap.String("key", k),
			zap.String("value", strings.Join(v, ",")),
		)
	}

	// Parse form data and log it
	r.ParseForm()
	for k, v := range r.Form {
		// Obfuscate sensitive data based on configured sensitive keys
		obfuscatedValues := internal.ObfuscateSensitiveData(k, v, m.SensitiveKeys)
		logger.Info("Form data",
			zap.String("key", k),
			zap.String("value", strings.Join(obfuscatedValues, ",")),
		)
	}

	// Log request details
	logger.Info("Request details",
		zap.String("method", r.Method),
		zap.String("protocol", r.Proto),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("request_uri", r.RequestURI),
		zap.String("user_agent", r.UserAgent()),
		zap.String("referer", r.Referer()),
	)

	// Call the next handler in the chain
	m.Next.ServeHTTP(w, r)

	// Calculate and log the request duration in milliseconds
	duration := time.Since(now).Milliseconds()
	logger.Info("Request duration",
		zap.Int64("duration_ms", duration),
	)
}
