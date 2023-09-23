package middleware

import (
	"net/http"
)

// CSPMiddleware is a middleware that sets the Content Security Policy (CSP) header.
type CSPMiddleware struct {
	Next http.Handler
	CSP  string // The CSP header value to set.
}

// NewCSPMiddleware creates a new CSPMiddleware instance.
func NewCSPMiddleware(next http.Handler, csp string) *CSPMiddleware {
	return &CSPMiddleware{
		Next: next,
		CSP:  csp,
	}
}

// ServeHTTP is the middleware handler function that sets the CSP header.
func (m *CSPMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", m.CSP)
	m.Next.ServeHTTP(w, r)
}
