package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// LogHandler handles log-related HTTP requests
type LogHandler struct {
	logger logging.Logger
}

// NewLogHandler creates a new LogHandler
func NewLogHandler(logger logging.Logger) *LogHandler {
	return &LogHandler{
		logger: logger,
	}
}

// CreateLog handles POST /logs
func (h *LogHandler) CreateLog(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode log request", "error", err)
		return WriteError(w, http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if req.Level == "" || req.Source == "" || req.Message == "" {
		return WriteError(w, http.StatusBadRequest, "level, source, and message are required")
	}

	// TODO: Implement log storage when LogUseCase is available
	// See issue Nexflow-6n4 for implementation details.
	// For now, just acknowledge the request
	h.logger.Info("log entry received", "level", req.Level, "source", req.Source, "message", req.Message)

	resp := &dto.LogResponse{
		Success: true,
		Log: &dto.LogDTO{
			Level:   req.Level,
			Source:  req.Source,
			Message: req.Message,
		},
	}

	return WriteJSON(w, http.StatusCreated, resp)
}

// ListLogs handles GET /logs
func (h *LogHandler) ListLogs(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement log listing when LogUseCase is available
	// See issue Nexflow-6n4 for implementation details.
	resp := &dto.LogsResponse{
		Success: true,
		Logs:    []*dto.LogDTO{},
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// RegisterLogRoutes registers log routes
func RegisterLogRoutes(r *Router, handler *LogHandler) {
	r.HandleFunc("POST /logs", handler.CreateLog)
	r.HandleFunc("GET /logs", handler.ListLogs)
}
