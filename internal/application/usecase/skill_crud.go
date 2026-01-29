package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

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
