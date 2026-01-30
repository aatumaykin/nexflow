package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// CreateSkill creates a new skill
func (uc *SkillUseCase) CreateSkill(ctx context.Context, req dto.CreateSkillRequest) (*dto.SkillResponse, error) {
	newSkill := entity.NewSkill(req.Name, req.Version, req.Location, req.Permissions, req.Metadata)

	if err := uc.skillRepo.Create(ctx, newSkill); err != nil {
		return handleSkillError(err, "failed to create skill")
	}

	uc.logger.Info("skill created", "skill_id", newSkill.ID, "name", newSkill.Name)

	return dto.SuccessSkillResponse(dto.SkillDTOFromEntity(newSkill)), nil
}
