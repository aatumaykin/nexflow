package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/usecase"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
	chatUseCase *usecase.ChatUseCase
	logger      Logger
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(chatUseCase *usecase.ChatUseCase, logger Logger) *TaskHandler {
	return &TaskHandler{
		chatUseCase: chatUseCase,
		logger:      logger,
	}
}

// GetSessionTasks handles GET /sessions/{id}/tasks
func (h *TaskHandler) GetSessionTasks(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		return WriteError(w, http.StatusBadRequest, "session id is required")
	}

	resp, err := h.chatUseCase.GetSessionTasks(ctx, sessionID)
	if err != nil {
		h.logger.Error("failed to get session tasks", "error", err, "session_id", sessionID)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// ExecuteSkill handles POST /skills/execute
func (h *TaskHandler) ExecuteSkill(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req dto.SkillExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode skill execution request", "error", err)
		return WriteError(w, http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if req.Skill == "" {
		return WriteError(w, http.StatusBadRequest, "skill name is required")
	}

	// Get session ID from query parameter
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		return WriteError(w, http.StatusBadRequest, "session_id is required")
	}

	resp, err := h.chatUseCase.ExecuteSkill(ctx, sessionID, req.Skill, req.Input)
	if err != nil {
		h.logger.Error("failed to execute skill", "error", err, "skill", req.Skill)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// RegisterTaskRoutes registers task routes
func RegisterTaskRoutes(r *Router, handler *TaskHandler) {
	r.HandleFunc("GET /sessions/{id}/tasks", handler.GetSessionTasks)
	r.HandleFunc("POST /skills/execute", handler.ExecuteSkill)
}
