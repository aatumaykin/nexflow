package http

import (
	"context"
	"net/http"
)

// Handler is an HTTP handler function that accepts a context and returns an error
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// ErrorHandler handles errors returned by Handler
type ErrorHandler interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error)
}

// DefaultErrorHandler is a simple error handler that returns JSON errors
type DefaultErrorHandler struct{}

// HandleError implements ErrorHandler interface
func (h *DefaultErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	WriteError(w, http.StatusInternalServerError, err.Error())
}

// HandlerAdapter adapts a Handler to an http.Handler
type HandlerAdapter struct {
	handler    Handler
	ctx        context.Context
	errHandler ErrorHandler
}

// NewHandlerAdapter creates a new handler adapter
func NewHandlerAdapter(handler Handler, ctx context.Context, errHandler ErrorHandler) http.Handler {
	if errHandler == nil {
		errHandler = &DefaultErrorHandler{}
	}

	return &HandlerAdapter{
		handler:    handler,
		ctx:        ctx,
		errHandler: errHandler,
	}
}

// ServeHTTP implements http.Handler interface
func (ha *HandlerAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := ha.handler(ha.ctx, w, r)
	if err != nil {
		ha.errHandler.HandleError(w, r, err)
	}
}

// HandlerBuilder helps build HTTP handlers with middleware
type HandlerBuilder struct {
	handler     http.Handler
	middlewares []Middleware
}

// NewHandlerBuilder creates a new handler builder
func NewHandlerBuilder(handler http.Handler) *HandlerBuilder {
	return &HandlerBuilder{
		handler:     handler,
		middlewares: make([]Middleware, 0),
	}
}

// Use adds a middleware to the chain
func (hb *HandlerBuilder) Use(middleware Middleware) *HandlerBuilder {
	hb.middlewares = append(hb.middlewares, middleware)
	return hb
}

// Build builds the final handler with all middleware applied
func (hb *HandlerBuilder) Build() http.Handler {
	handler := hb.handler

	// Apply middlewares in reverse order
	for i := len(hb.middlewares) - 1; i >= 0; i-- {
		handler = hb.middlewares[i](handler)
	}

	return handler
}
