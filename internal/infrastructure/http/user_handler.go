package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/usecase"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase *usecase.UserUseCase
	logger      Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase *usecase.UserUseCase, logger Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode user request", "error", err)
		return WriteError(w, http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if req.Channel == "" || req.ChannelID == "" {
		return WriteError(w, http.StatusBadRequest, "channel and channel_id are required")
	}

	resp, err := h.userUseCase.CreateUser(ctx, req)
	if err != nil {
		h.logger.Error("failed to create user", "error", err)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusCreated, resp)
}

// GetUserByID handles GET /users/{id}
func (h *UserHandler) GetUserByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return WriteError(w, http.StatusBadRequest, "user id is required")
	}

	resp, err := h.userUseCase.GetUserByID(ctx, id)
	if err != nil {
		h.logger.Error("failed to get user", "error", err, "user_id", id)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusNotFound, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// GetUserByChannel handles GET /users/channel/{channel}/{channelID}
func (h *UserHandler) GetUserByChannel(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	channel := r.PathValue("channel")
	channelID := r.PathValue("channelID")

	if channel == "" || channelID == "" {
		return WriteError(w, http.StatusBadRequest, "channel and channel_id are required")
	}

	resp, err := h.userUseCase.GetUserByChannel(ctx, channel, channelID)
	if err != nil {
		h.logger.Error("failed to get user by channel", "error", err, "channel", channel, "channel_id", channelID)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusNotFound, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	resp, err := h.userUseCase.ListUsers(ctx)
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// DeleteUser handles DELETE /users/{id}
func (h *UserHandler) DeleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return WriteError(w, http.StatusBadRequest, "user id is required")
	}

	resp, err := h.userUseCase.DeleteUser(ctx, id)
	if err != nil {
		h.logger.Error("failed to delete user", "error", err, "user_id", id)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusNotFound, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// RegisterUserRoutes registers user routes
func RegisterUserRoutes(r *Router, handler *UserHandler) {
	r.HandleFunc("POST /users", handler.CreateUser)
	r.HandleFunc("GET /users", handler.ListUsers)
	r.HandleFunc("GET /users/{id}", handler.GetUserByID)
	r.HandleFunc("GET /users/channel/{channel}/{channelID}", handler.GetUserByChannel)
	r.HandleFunc("DELETE /users/{id}", handler.DeleteUser)
}

// Logger interface for structured logging
type Logger interface {
	Error(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}
