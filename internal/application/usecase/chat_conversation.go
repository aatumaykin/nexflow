package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// GetConversation retrieves conversation history for a session
func (uc *ChatUseCase) GetConversation(ctx context.Context, sessionID string) (*dto.MessagesResponse, error) {
	messages, err := uc.messageRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return dto.ErrorMessageResponse(fmt.Errorf("failed to get conversation: %w", err)), fmt.Errorf("failed to get conversation: %w", err)
	}

	messageDTOs := make([]*dto.MessageDTO, 0, len(messages))
	for _, msg := range messages {
		messageDTOs = append(messageDTOs, dto.MessageDTOFromEntity(msg))
	}

	return dto.SuccessMessagesResponse(messageDTOs), nil
}
