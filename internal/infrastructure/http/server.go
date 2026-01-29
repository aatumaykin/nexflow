package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Server represents HTTP server
type Server struct {
	httpServer  *http.Server
	middlewares []Middleware
}

// ServerConfig holds configuration for HTTP server
type ServerConfig struct {
	Addr         string // Server address (e.g., ":8080")
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Handler      http.Handler
}

// NewServer creates a new HTTP server
func NewServer(cfg *ServerConfig) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Addr,
			Handler:      cfg.Handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		middlewares: make([]Middleware, 0),
	}
}

// Start starts HTTP server
func (s *Server) Start(ctx context.Context) error {
	log.Printf("Starting HTTP server on %s", s.httpServer.Addr)

	errChan := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start HTTP server: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return s.Shutdown(ctx)
	}
}

// Shutdown gracefully shuts down HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	log.Println("HTTP server shutdown complete")
	return nil
}

// Use adds a middleware to the server
func (s *Server) Use(middleware Middleware) {
	s.middlewares = append(s.middlewares, middleware)
}

// Addr returns the server address
func (s *Server) Addr() string {
	return s.httpServer.Addr
}
