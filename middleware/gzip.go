package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipMiddleware is a middleware that compresses response data using gzip.
type GzipMiddleware struct {
	Next             http.Handler
	CompressionLevel *int // Compression level (0-9), where 0 is no compression, and 9 is maximum compression.
}

// NewGzipMiddleware creates a new GzipMiddleware instance with the specified compression level.
// If compressionLevel is nil, the default compression level (-1) is used.
func NewGzipMiddleware(next http.Handler, compressionLevel *int) *GzipMiddleware {
	return &GzipMiddleware{
		Next:             next,
		CompressionLevel: compressionLevel,
	}
}

// ServeHTTP is the middleware handler function that performs gzip compression.
func (m *GzipMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if the client supports gzip encoding.
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		var compressionLevel int
		if m.CompressionLevel != nil {
			compressionLevel = *m.CompressionLevel
		} else {
			// Use the default compression level (-1) if not specified.
			compressionLevel = -1
		}

		// Create a gzip.Writer with the specified compression level.
		gz, err := gzip.NewWriterLevel(w, compressionLevel)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		// Set the Content-Encoding header to indicate gzip encoding.
		w.Header().Set("Content-Encoding", "gzip")

		// Wrap the response writer with the gzip writer.
		m.Next.ServeHTTP(NewGzipResponseWriter(w, gz), r)
	} else {
		// If the client does not support gzip, pass the request along as-is.
		m.Next.ServeHTTP(w, r)
	}
}

// GzipResponseWriter is a custom response writer that wraps a gzip writer.
type GzipResponseWriter struct {
	http.ResponseWriter
	gz *gzip.Writer
}

// NewGzipResponseWriter creates a new GzipResponseWriter instance.
func NewGzipResponseWriter(w http.ResponseWriter, gz *gzip.Writer) *GzipResponseWriter {
	return &GzipResponseWriter{
		ResponseWriter: w,
		gz:             gz,
	}
}

// Write writes compressed data to the gzip writer.
func (grw *GzipResponseWriter) Write(b []byte) (int, error) {
	return grw.gz.Write(b)
}
