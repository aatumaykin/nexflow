package http

import (
	"net/http"
)

// Router wraps http.ServeMux with handler registration
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new Router
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// HandleFunc registers a handler for the given pattern
func (r *Router) HandleFunc(pattern string, handler Handler) {
	r.mux.Handle(pattern, NewHandlerAdapter(handler, nil, nil))
}

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Handler returns the underlying http.Handler
func (r *Router) Handler() http.Handler {
	return r.mux
}
