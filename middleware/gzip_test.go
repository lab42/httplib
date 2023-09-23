package middleware_test

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lab42/httplib/middleware"
)

func TestGzipMiddleware(t *testing.T) {
	// Create a sample HTTP handler for testing.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Create a test request.
	req := httptest.NewRequest("GET", "/", nil)

	// Set an Accept-Encoding header to indicate gzip support.
	req.Header.Set("Accept-Encoding", "gzip")

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Create a GzipMiddleware instance.
	middleware := middleware.NewGzipMiddleware(handler, nil)

	// Execute the middleware.
	middleware.ServeHTTP(rr, req)

	// Check if the response has gzip encoding.
	if rr.Header().Get("Content-Encoding") != "gzip" {
		t.Errorf("Expected Content-Encoding to be 'gzip', but got '%s'", rr.Header().Get("Content-Encoding"))
	}

	// Decode the gzip content to verify the response.
	gr, err := gzip.NewReader(rr.Body)
	if err != nil {
		t.Errorf("Failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	var decompressed bytes.Buffer
	_, err = decompressed.ReadFrom(gr)
	if err != nil {
		t.Errorf("Failed to read gzip content: %v", err)
	}

	expectedResponse := "Hello, World!"
	if decompressed.String() != expectedResponse {
		t.Errorf("Expected response '%s', but got '%s'", expectedResponse, decompressed.String())
	}
}

func TestGzipResponseWriter_Write(t *testing.T) {
	// Create a mock response writer.
	mockResponseWriter := httptest.NewRecorder()

	// Create a gzip.Writer for testing.
	gz, err := gzip.NewWriterLevel(mockResponseWriter, gzip.DefaultCompression)
	if err != nil {
		t.Fatalf("Failed to create gzip writer: %v", err)
	}
	defer gz.Close()

	// Create a GzipResponseWriter instance.
	grw := middleware.NewGzipResponseWriter(mockResponseWriter, gz)

	// Write some data using the GzipResponseWriter.
	data := []byte("Hello, Gzip!")
	_, err = grw.Write(data)
	if err != nil {
		t.Errorf("Failed to write to GzipResponseWriter: %v", err)
	}

	// Check if the data was written to the gzip.Writer.
	gz.Flush()
	if mockResponseWriter.Body.String() == "" {
		t.Error("Expected data to be written to the response writer, but got an empty response")
	}
}
