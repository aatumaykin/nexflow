package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
)

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

	// Update skill fields
	if err := uc.updateSkillFields(skill, req); err != nil {
		return &dto.SkillResponse{
			Success: false,
			Error:   err.Error(),
		}, err
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

// updateSkillFields updates skill fields from request
func (uc *SkillUseCase) updateSkillFields(skill *entity.Skill, req dto.UpdateSkillRequest) error {
	// Update version
	if req.Version != "" {
		skill.Version = valueobject.MustNewVersion(req.Version)
	}

	// Update location
	if req.Location != "" {
		skill.Location = req.Location
	}

	// Update permissions
	if req.Permissions != nil {
		permissionsJSON, err := dto.SliceToString(req.Permissions)
		if err != nil {
			return fmt.Errorf("failed to marshal permissions: %w", err)
		}
		skill.Permissions = permissionsJSON
	}

	// Update metadata
	if req.Metadata != nil {
		metadataJSON, err := dto.MapToString(req.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		skill.Metadata = metadataJSON
	}

	return nil
}
