package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// SkillUseCase handles skill-related business logic
type SkillUseCase struct {
	skillRepo    repository.SkillRepository
	skillRuntime ports.SkillRuntime
	logger       logging.Logger
}

// NewSkillUseCase creates a new SkillUseCase
func NewSkillUseCase(
	skillRepo repository.SkillRepository,
	skillRuntime ports.SkillRuntime,
	logger logging.Logger,
) *SkillUseCase {
	return &SkillUseCase{
		skillRepo:    skillRepo,
		skillRuntime: skillRuntime,
		logger:       logger,
	}
}

// CreateSkill creates a new skill
func (uc *SkillUseCase) CreateSkill(ctx context.Context, req dto.CreateSkillRequest) (*dto.SkillResponse, error) {
	// Create new skill using entity constructor
	newSkill := entity.NewSkill(req.Name, req.Version, req.Location, req.Permissions, req.Metadata)

	if err := uc.skillRepo.Create(ctx, newSkill); err != nil {
		return &dto.SkillResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create skill: %v", err),
		}, fmt.Errorf("failed to create skill: %w", err)
	}

	uc.logger.Info("skill created", "skill_id", newSkill.ID, "name", newSkill.Name)

	return &dto.SkillResponse{
		Success: true,
		Skill:   dto.SkillDTOFromEntity(newSkill),
	}, nil
}

// GetSkillByID retrieves a skill by ID
func (uc *SkillUseCase) GetSkillByID(ctx context.Context, id string) (*dto.SkillResponse, error) {
	skill, err := uc.skillRepo.FindByID(ctx, id)
	if err != nil {
		return &dto.SkillResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to find skill: %v", err),
		}, fmt.Errorf("failed to find skill: %w", err)
	}

	return &dto.SkillResponse{
		Success: true,
		Skill:   dto.SkillDTOFromEntity(skill),
	}, nil
}

// GetSkillByName retrieves a skill by name
func (uc *SkillUseCase) GetSkillByName(ctx context.Context, name string) (*dto.SkillResponse, error) {
	skill, err := uc.skillRepo.FindByName(ctx, name)
	if err != nil {
		return &dto.SkillResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to find skill by name: %v", err),
		}, fmt.Errorf("failed to find skill by name: %w", err)
	}

	return &dto.SkillResponse{
		Success: true,
		Skill:   dto.SkillDTOFromEntity(skill),
	}, nil
}

// ListSkills retrieves all skills
func (uc *SkillUseCase) ListSkills(ctx context.Context) (*dto.SkillsResponse, error) {
	skills, err := uc.skillRepo.List(ctx)
	if err != nil {
		return &dto.SkillsResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to list skills: %v", err),
		}, fmt.Errorf("failed to list skills: %w", err)
	}

	skillDTOs := make([]*dto.SkillDTO, 0, len(skills))
	for _, skill := range skills {
		skillDTOs = append(skillDTOs, dto.SkillDTOFromEntity(skill))
	}

	return &dto.SkillsResponse{
		Success: true,
		Skills:  skillDTOs,
	}, nil
}

// UpdateSkill updates an existing skill
func (uc *SkillUseCase) UpdateSkill(ctx context.Context, id string, req dto.UpdateSkillRequest) (*dto.SkillResponse, error) {
	// Get existing skill
	skill, err := uc.skillRepo.FindByID(ctx, id)
	if err != nil {
		return &dto.SkillResponse{
			Success: false,
			Error:   fmt.Sprintf("skill not found: %v", err),
		}, fmt.Errorf("skill not found: %w", err)
	}

	// Update fields
	if req.Version != "" {
		skill.Version = req.Version
	}
	if req.Location != "" {
		skill.Location = req.Location
	}
	if req.Permissions != nil {
		permissionsJSON, err := dto.SliceToString(req.Permissions)
		if err != nil {
			return &dto.SkillResponse{
				Success: false,
				Error:   fmt.Sprintf("failed to marshal permissions: %v", err),
			}, fmt.Errorf("failed to marshal permissions: %w", err)
		}
		skill.Permissions = permissionsJSON
	}
	if req.Metadata != nil {
		metadataJSON, err := dto.MapToString(req.Metadata)
		if err != nil {
			return &dto.SkillResponse{
				Success: false,
				Error:   fmt.Sprintf("failed to marshal metadata: %v", err),
			}, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		skill.Metadata = metadataJSON
	}

	// Save updated skill
	if err := uc.skillRepo.Update(ctx, skill); err != nil {
		return &dto.SkillResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to update skill: %v", err),
		}, fmt.Errorf("failed to update skill: %w", err)
	}

	uc.logger.Info("skill updated", "skill_id", skill.ID, "name", skill.Name)

	return &dto.SkillResponse{
		Success: true,
		Skill:   dto.SkillDTOFromEntity(skill),
	}, nil
}

// DeleteSkill deletes a skill by ID
func (uc *SkillUseCase) DeleteSkill(ctx context.Context, id string) (*dto.SkillResponse, error) {
	// Check if skill exists
	skill, err := uc.skillRepo.FindByID(ctx, id)
	if err != nil {
		return &dto.SkillResponse{
			Success: false,
			Error:   fmt.Sprintf("skill not found: %v", err),
		}, fmt.Errorf("skill not found: %w", err)
	}

	// Delete skill
	if err := uc.skillRepo.Delete(ctx, id); err != nil {
		return &dto.SkillResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to delete skill: %v", err),
		}, fmt.Errorf("failed to delete skill: %w", err)
	}

	uc.logger.Info("skill deleted", "skill_id", skill.ID, "name", skill.Name)

	return &dto.SkillResponse{
		Success: true,
		Skill:   dto.SkillDTOFromEntity(skill),
	}, nil
}

// ExecuteSkill executes a skill with given input parameters
func (uc *SkillUseCase) ExecuteSkill(ctx context.Context, req dto.SkillExecutionRequest) (*dto.SkillExecutionResponse, error) {
	// Validate skill exists
	_, err := uc.skillRepo.FindByName(ctx, req.Skill)
	if err != nil {
		return &dto.SkillExecutionResponse{
			Success: false,
			Error:   fmt.Sprintf("skill not found: %v", err),
		}, fmt.Errorf("skill not found: %w", err)
	}

	// Execute skill through runtime
	result, err := uc.skillRuntime.Execute(ctx, req.Skill, req.Input)
	if err != nil {
		return &dto.SkillExecutionResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to execute skill: %v", err),
		}, fmt.Errorf("failed to execute skill: %w", err)
	}

	uc.logger.Info("skill executed", "skill", req.Skill, "success", result.Success)

	return &dto.SkillExecutionResponse{
		Success: result.Success,
		Output:  result.Output,
		Error:   result.Error,
	}, nil
}

// ValidateSkill validates a skill
func (uc *SkillUseCase) ValidateSkill(ctx context.Context, skillName string) error {
	return uc.skillRuntime.Validate(skillName)
}

// ListAvailableSkills returns list of available skill names
func (uc *SkillUseCase) ListAvailableSkills(ctx context.Context) ([]string, error) {
	return uc.skillRuntime.List()
}

// GetSkillDetails returns detailed skill information
func (uc *SkillUseCase) GetSkillDetails(ctx context.Context, skillName string) (map[string]interface{}, error) {
	return uc.skillRuntime.GetSkill(skillName)
}
