package middleware

import (
	"net/http"
)

// CSRFMiddleware is a middleware that provides Cross-Site Request Forgery (CSRF) protection.
type CSRFMiddleware struct {
	Next        http.Handler
	CSRFHeader  string // The header containing the CSRF token.
	CSRFCookie  string // The name of the CSRF token cookie.
	CSRFParam   string // The name of the CSRF token parameter in form submissions.
	CSRFToken   string // The expected CSRF token value.
	ErrorStatus int    // The HTTP status code to use when CSRF validation fails (e.g., http.StatusForbidden).
}

// NewCSRFMiddleware creates a new CSRFMiddleware instance.
func NewCSRFMiddleware(next http.Handler, csrfHeader, csrfCookie, csrfParam, csrfToken string, errorStatus int) *CSRFMiddleware {
	return &CSRFMiddleware{
		Next:        next,
		CSRFHeader:  csrfHeader,
		CSRFCookie:  csrfCookie,
		CSRFParam:   csrfParam,
		CSRFToken:   csrfToken,
		ErrorStatus: errorStatus,
	}
}

// ServeHTTP is the middleware handler function that provides CSRF protection.
func (m *CSRFMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Retrieve the CSRF token from the request header, cookie, or form parameter.
	token := r.Header.Get(m.CSRFHeader)
	if token == "" {
		cookie, err := r.Cookie(m.CSRFCookie)
		if err == nil {
			token = cookie.Value
		}
	}
	if token == "" {
		token = r.FormValue(m.CSRFParam)
	}

	// Check if the token matches the expected value.
	if token != m.CSRFToken {
		http.Error(w, "CSRF token validation failed", m.ErrorStatus)
		return
	}

	m.Next.ServeHTTP(w, r)
}
