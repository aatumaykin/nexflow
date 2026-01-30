package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// DeleteSkill deletes a skill by ID
func (uc *SkillUseCase) DeleteSkill(ctx context.Context, id string) (*dto.SkillResponse, error) {
	skill, err := uc.skillRepo.FindByID(ctx, id)
	if err != nil {
		return handleSkillError(err, "skill not found")
	}

	if err := uc.skillRepo.Delete(ctx, id); err != nil {
		return handleSkillError(err, "failed to delete skill")
	}

	uc.logger.Info("skill deleted", "skill_id", skill.ID, "name", skill.Name)

	return dto.SuccessSkillResponse(dto.SkillDTOFromEntity(skill)), nil
}
