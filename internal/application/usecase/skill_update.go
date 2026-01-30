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
	skill, err := uc.skillRepo.FindByID(ctx, id)
	if err != nil {
		return handleSkillError(err, "skill not found")
	}

	if err := uc.updateSkillFields(skill, req); err != nil {
		return handleSkillError(err, "failed to update skill fields")
	}

	if err := uc.skillRepo.Update(ctx, skill); err != nil {
		return handleSkillError(err, "failed to update skill")
	}

	uc.logger.Info("skill updated", "skill_id", skill.ID, "name", skill.Name)

	return dto.SuccessSkillResponse(dto.SkillDTOFromEntity(skill)), nil
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
