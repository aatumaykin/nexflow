package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

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
