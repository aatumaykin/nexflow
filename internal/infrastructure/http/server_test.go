package http

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cfg := &ServerConfig{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	server := NewServer(cfg)

	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.Addr())
}

func TestServer_Use(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cfg := &ServerConfig{
		Addr:    ":8080",
		Handler: handler,
	}

	server := NewServer(cfg)

	server.Use(Recovery)

	assert.Len(t, server.middlewares, 1)
}

func TestServer_Start(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	cfg := &ServerConfig{
		Addr:    ":18080", // Use different port to avoid conflicts
		Handler: handler,
	}

	server := NewServer(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := server.Start(ctx)

	// Server should start without error (context will cancel it)
	// In real scenario, you'd wait for server to be ready
	assert.NoError(t, err)
}

func TestServer_Shutdown(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cfg := &ServerConfig{
		Addr:    ":18081", // Use different port to avoid conflicts
		Handler: handler,
	}

	server := NewServer(cfg)

	// Start server in background
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = server.Start(ctx)
	}()

	// Wait a bit for server to start
	time.Sleep(50 * time.Millisecond)

	// Shutdown server
	err := server.Shutdown(context.Background())

	require.NoError(t, err)

	// Cancel start context
	cancel()
}
