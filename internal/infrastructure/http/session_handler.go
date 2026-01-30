package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/usecase"
)

// SessionHandler handles session-related HTTP requests
type SessionHandler struct {
	chatUseCase *usecase.ChatUseCase
	logger      Logger
}

// NewSessionHandler creates a new SessionHandler
func NewSessionHandler(chatUseCase *usecase.ChatUseCase, logger Logger) *SessionHandler {
	return &SessionHandler{
		chatUseCase: chatUseCase,
		logger:      logger,
	}
}

// CreateSession handles POST /sessions
func (h *SessionHandler) CreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode session request", "error", err)
		return WriteError(w, http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if req.UserID == "" {
		return WriteError(w, http.StatusBadRequest, "user_id is required")
	}

	resp, err := h.chatUseCase.CreateSession(ctx, req)
	if err != nil {
		h.logger.Error("failed to create session", "error", err)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusCreated, resp)
}

// GetUserSessions handles GET /users/{id}/sessions
func (h *SessionHandler) GetUserSessions(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userID := r.PathValue("id")
	if userID == "" {
		return WriteError(w, http.StatusBadRequest, "user id is required")
	}

	resp, err := h.chatUseCase.GetUserSessions(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get user sessions", "error", err, "user_id", userID)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// RegisterSessionRoutes registers session routes
func RegisterSessionRoutes(r *Router, handler *SessionHandler) {
	r.HandleFunc("POST /sessions", handler.CreateSession)
	r.HandleFunc("GET /users/{id}/sessions", handler.GetUserSessions)
}
