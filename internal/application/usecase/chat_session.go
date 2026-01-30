package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// GetUserSessions retrieves all sessions for a user
func (uc *ChatUseCase) GetUserSessions(ctx context.Context, userID string) (*dto.SessionsResponse, error) {
	sessions, err := uc.sessionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return dto.ErrorSessionsResponse(err), err
	}

	sessionDTOs := make([]*dto.SessionDTO, 0, len(sessions))
	for _, session := range sessions {
		sessionDTOs = append(sessionDTOs, dto.SessionDTOFromEntity(session))
	}

	return dto.SuccessSessionsResponse(sessionDTOs), nil
}

// CreateSession creates a new session for a user
func (uc *ChatUseCase) CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	session := entity.NewSession(req.UserID)
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return handleSessionError(err, "failed to create session")
	}

	return dto.SuccessSessionResponse(dto.SessionDTOFromEntity(session)), nil
}
