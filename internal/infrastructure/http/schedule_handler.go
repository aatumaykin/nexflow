package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/usecase"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// ScheduleHandler handles schedule-related HTTP requests
type ScheduleHandler struct {
	scheduleUseCase *usecase.ScheduleUseCase
	logger          logging.Logger
}

// NewScheduleHandler creates a new ScheduleHandler
func NewScheduleHandler(scheduleUseCase *usecase.ScheduleUseCase, logger logging.Logger) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleUseCase: scheduleUseCase,
		logger:          logger,
	}
}

// CreateSchedule handles POST /schedules
func (h *ScheduleHandler) CreateSchedule(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode schedule request", "error", err)
		return WriteError(w, http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if req.Skill == "" || req.CronExpression == "" {
		return WriteError(w, http.StatusBadRequest, "skill and cron_expression are required")
	}

	resp, err := h.scheduleUseCase.CreateSchedule(ctx, req)
	if err != nil {
		h.logger.Error("failed to create schedule", "error", err)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusCreated, resp)
}

// GetScheduleByID handles GET /schedules/{id}
func (h *ScheduleHandler) GetScheduleByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return WriteError(w, http.StatusBadRequest, "schedule id is required")
	}

	resp, err := h.scheduleUseCase.GetScheduleByID(ctx, id)
	if err != nil {
		h.logger.Error("failed to get schedule", "error", err, "schedule_id", id)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusNotFound, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// ListSchedules handles GET /schedules
func (h *ScheduleHandler) ListSchedules(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	enabled := r.URL.Query().Get("enabled")
	var resp *dto.SchedulesResponse
	var err error

	if enabled == "true" {
		resp, err = h.scheduleUseCase.ListEnabledSchedules(ctx)
	} else {
		resp, err = h.scheduleUseCase.ListSchedules(ctx)
	}

	if err != nil {
		h.logger.Error("failed to list schedules", "error", err)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// GetSchedulesBySkill handles GET /skills/{skill}/schedules
func (h *ScheduleHandler) GetSchedulesBySkill(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	skill := r.PathValue("skill")
	if skill == "" {
		return WriteError(w, http.StatusBadRequest, "skill name is required")
	}

	resp, err := h.scheduleUseCase.GetSchedulesBySkill(ctx, skill)
	if err != nil {
		h.logger.Error("failed to get schedules by skill", "error", err, "skill", skill)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// UpdateSchedule handles PUT /schedules/{id}
func (h *ScheduleHandler) UpdateSchedule(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return WriteError(w, http.StatusBadRequest, "schedule id is required")
	}

	var req dto.UpdateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode schedule update request", "error", err)
		return WriteError(w, http.StatusBadRequest, "invalid request body")
	}

	resp, err := h.scheduleUseCase.UpdateSchedule(ctx, id, req)
	if err != nil {
		h.logger.Error("failed to update schedule", "error", err, "schedule_id", id)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// ToggleSchedule handles POST /schedules/{id}/toggle
func (h *ScheduleHandler) ToggleSchedule(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return WriteError(w, http.StatusBadRequest, "schedule id is required")
	}

	var req dto.ToggleScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode toggle schedule request", "error", err)
		return WriteError(w, http.StatusBadRequest, "invalid request body")
	}

	resp, err := h.scheduleUseCase.ToggleSchedule(ctx, id, req)
	if err != nil {
		h.logger.Error("failed to toggle schedule", "error", err, "schedule_id", id)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// EnableSchedule handles POST /schedules/{id}/enable
func (h *ScheduleHandler) EnableSchedule(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return WriteError(w, http.StatusBadRequest, "schedule id is required")
	}

	resp, err := h.scheduleUseCase.EnableSchedule(ctx, id)
	if err != nil {
		h.logger.Error("failed to enable schedule", "error", err, "schedule_id", id)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// DisableSchedule handles POST /schedules/{id}/disable
func (h *ScheduleHandler) DisableSchedule(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	if id == "" {
		return WriteError(w, http.StatusBadRequest, "schedule id is required")
	}

	resp, err := h.scheduleUseCase.DisableSchedule(ctx, id)
	if err != nil {
		h.logger.Error("failed to disable schedule", "error", err, "schedule_id", id)
		return WriteError(w, http.StatusInternalServerError, resp.Error)
	}

	if !resp.Success {
		return WriteError(w, http.StatusBadRequest, resp.Error)
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// RegisterScheduleRoutes registers schedule routes
func RegisterScheduleRoutes(r *Router, handler *ScheduleHandler) {
	r.HandleFunc("POST /schedules", handler.CreateSchedule)
	r.HandleFunc("GET /schedules", handler.ListSchedules)
	r.HandleFunc("GET /schedules/{id}", handler.GetScheduleByID)
	r.HandleFunc("GET /skills/{skill}/schedules", handler.GetSchedulesBySkill)
	r.HandleFunc("PUT /schedules/{id}", handler.UpdateSchedule)
	r.HandleFunc("POST /schedules/{id}/toggle", handler.ToggleSchedule)
	r.HandleFunc("POST /schedules/{id}/enable", handler.EnableSchedule)
	r.HandleFunc("POST /schedules/{id}/disable", handler.DisableSchedule)
}
