package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestLogRequestMiddleware_Success(t *testing.T) {
	// Replace zap global logger with a test logger
	logger := zaptest.NewLogger(t)
	orig := zap.L()
	zap.ReplaceGlobals(logger)
	defer zap.ReplaceGlobals(orig)

	// Prepare a request with a body
	body := "test body"
	req := httptest.NewRequest(http.MethodPost, "/test?foo=bar", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	// Dummy handler to check if next is called
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		// Body should still be readable
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("unexpected error reading body: %v", err)
		}
		if string(b) != body {
			t.Errorf("expected body %q, got %q", body, string(b))
		}
	})

	mw := LogRequestMiddleware()
	mw(next).ServeHTTP(rec, req)

	if !called {
		t.Error("expected next handler to be called")
	}
}

func TestLogRequestMiddleware_PeekBodyError(t *testing.T) {
	// Replace zap global logger with a test logger
	logger := zaptest.NewLogger(t)
	orig := zap.L()
	zap.ReplaceGlobals(logger)
	defer zap.ReplaceGlobals(orig)

	// Create a request with a body that errors on Read
	req := httptest.NewRequest(http.MethodPost, "/fail", errReader{})
	rec := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called on error")
	})

	mw := LogRequestMiddleware()
	mw(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
