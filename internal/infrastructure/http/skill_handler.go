package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/usecase"
)

// SkillHandler handles skill-related HTTP requests
type SkillHandler struct {
	skillUseCase *usecase.SkillUseCase
	logger       Logger
}

// NewSkillHandler creates a new SkillHandler
func NewSkillHandler(skillUseCase *usecase.SkillUseCase, logger Logger) *SkillHandler {
	return &SkillHandler{
		skillUseCase: skillUseCase,
		logger:       logger,
	}
}

// CreateSkill handles POST /skills
func (h *SkillHandler) CreateSkill(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateSkillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode skill request", "error", err)
		return WriteError(w, http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if req.Name == "" || req.Version == "" || req.Location == "" {
		return WriteError(w, http.StatusBadRequest, "name, version, and location are required")
	}

	resp, err := h.skillUseCase.CreateSkill(ctx, req)
	if err != nil {
		h.logger.Error("failed to create skill", "error", err)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusCreated, resp)
}

// GetSkillByID handles GET /skills/{id}
func (h *SkillHandler) GetSkillByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return WriteError(w, http.StatusBadRequest, "skill id is required")
	}

	resp, err := h.skillUseCase.GetSkillByID(ctx, id)
	if err != nil {
		h.logger.Error("failed to get skill", "error", err, "skill_id", id)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusNotFound, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// GetSkillByName handles GET /skills/name/{name}
func (h *SkillHandler) GetSkillByName(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	name := r.PathValue("name")
	if name == "" {
		return WriteError(w, http.StatusBadRequest, "skill name is required")
	}

	resp, err := h.skillUseCase.GetSkillByName(ctx, name)
	if err != nil {
		h.logger.Error("failed to get skill by name", "error", err, "skill_name", name)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusNotFound, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// ListSkills handles GET /skills
func (h *SkillHandler) ListSkills(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	resp, err := h.skillUseCase.ListSkills(ctx)
	if err != nil {
		h.logger.Error("failed to list skills", "error", err)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// RegisterSkillRoutes registers skill routes
func RegisterSkillRoutes(r *Router, handler *SkillHandler) {
	r.HandleFunc("POST /skills", handler.CreateSkill)
	r.HandleFunc("GET /skills", handler.ListSkills)
	r.HandleFunc("GET /skills/{id}", handler.GetSkillByID)
	r.HandleFunc("GET /skills/name/{name}", handler.GetSkillByName)
}
