package logging

import (
	"net/http"
	"time"
)

// Middleware creates a logging middleware that logs HTTP requests.
// By default, it uses NoopLogger that does nothing.
// Use WithMiddlewareLogger option to provide a custom logger.
func Middleware(opts ...MiddlewareOption) func(http.Handler) http.Handler {
	// Apply options with default NoopLogger
	config := &middlewareConfig{
		logger: NewNoopLogger(),
	}
	for _, opt := range opts {
		opt(config)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w}

			// Log request
			config.logger.Info("Request started",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)

			// Call next handler
			next.ServeHTTP(wrapped, r)

			// Log response
			duration := time.Since(start)
			config.logger.Info("Request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.status,
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}

// MiddlewareOption is a function that configures the logging middleware.
type MiddlewareOption func(*middlewareConfig)

// WithMiddlewareLogger sets the logger for the middleware.
func WithMiddlewareLogger(logger Logger) MiddlewareOption {
	return func(cfg *middlewareConfig) {
		cfg.logger = logger
	}
}

// middlewareConfig holds configuration for the logging middleware.
type middlewareConfig struct {
	logger Logger
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
