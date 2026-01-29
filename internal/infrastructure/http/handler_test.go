package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockErrorHandler for testing
type MockErrorHandler struct {
	called     bool
	lastError  string
	statusCode int
}

func (m *MockErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	m.called = true
	m.lastError = err.Error()
	w.WriteHeader(http.StatusInternalServerError)
}

func TestHandlerAdapter_ServeHTTP_Success(t *testing.T) {
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	adapter := NewHandlerAdapter(handler, context.Background(), nil)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	adapter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerAdapter_ServeHTTP_Error(t *testing.T) {
	mockErr := assert.AnError
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return mockErr
	}

	mockErrorHandler := &MockErrorHandler{statusCode: http.StatusInternalServerError}
	adapter := NewHandlerAdapter(handler, context.Background(), mockErrorHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	adapter.ServeHTTP(w, req)

	assert.True(t, mockErrorHandler.called)
	assert.Equal(t, mockErr.Error(), mockErrorHandler.lastError)
}

func TestNewHandlerAdapter_DefaultErrorHandler(t *testing.T) {
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return assert.AnError
	}

	adapter := NewHandlerAdapter(handler, context.Background(), nil)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	adapter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandlerBuilder_Use(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	builder := NewHandlerBuilder(handler)
	middlewareCalled := false

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	builder.Use(middleware)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	builder.Build().ServeHTTP(w, req)

	assert.True(t, middlewareCalled)
}

func TestHandlerBuilder_Build(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	builder := NewHandlerBuilder(handler)
	middleware1Called := false
	middleware2Called := false

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware1Called = true
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware2Called = true
			next.ServeHTTP(w, r)
		})
	}

	builder.Use(middleware1)
	builder.Use(middleware2)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	builder.Build().ServeHTTP(w, req)

	// Both middlewares should be called
	assert.True(t, middleware1Called)
	assert.True(t, middleware2Called)
}
