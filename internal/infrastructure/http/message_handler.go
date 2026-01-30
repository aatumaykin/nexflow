package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/usecase"
)

// MessageHandler handles message-related HTTP requests
type MessageHandler struct {
	chatUseCase *usecase.ChatUseCase
	logger      Logger
}

// NewMessageHandler creates a new MessageHandler
func NewMessageHandler(chatUseCase *usecase.ChatUseCase, logger Logger) *MessageHandler {
	return &MessageHandler{
		chatUseCase: chatUseCase,
		logger:      logger,
	}
}

// GetConversation handles GET /sessions/{id}/messages
func (h *MessageHandler) GetConversation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		return WriteError(w, http.StatusBadRequest, "session id is required")
	}

	resp, err := h.chatUseCase.GetConversation(ctx, sessionID)
	if err != nil {
		h.logger.Error("failed to get conversation", "error", err, "session_id", sessionID)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// SendMessage handles POST /chat/send
func (h *MessageHandler) SendMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode send message request", "error", err)
		return WriteError(w, http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if req.UserID == "" {
		return WriteError(w, http.StatusBadRequest, "user_id is required")
	}
	if req.Message.Content == "" {
		return WriteError(w, http.StatusBadRequest, "message content is required")
	}

	resp, err := h.chatUseCase.SendMessage(ctx, req)
	if err != nil {
		h.logger.Error("failed to send message", "error", err)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// RegisterMessageRoutes registers message routes
func RegisterMessageRoutes(r *Router, handler *MessageHandler) {
	r.HandleFunc("GET /sessions/{id}/messages", handler.GetConversation)
	r.HandleFunc("POST /chat/send", handler.SendMessage)
}
