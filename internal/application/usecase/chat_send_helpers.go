package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// findOrCreateUser finds existing user or creates new one
func (uc *ChatUseCase) findOrCreateUser(ctx context.Context, userID string) (*entity.User, error) {
	user, err := uc.userRepo.FindByChannel(ctx, "web", userID)
	if err != nil {
		newUser := entity.NewUser("web", userID)
		if err := uc.userRepo.Create(ctx, newUser); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		return newUser, nil
	}
	return user, nil
}

// createSession creates a new session for the user
func (uc *ChatUseCase) createSession(ctx context.Context, user *entity.User) (*entity.Session, error) {
	session := entity.NewSession(string(user.ID))
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return session, nil
}

// saveUserMessage saves user message to repository
func (uc *ChatUseCase) saveUserMessage(ctx context.Context, session *entity.Session, content string) (*entity.Message, error) {
	userMessage := entity.NewUserMessage(string(session.ID), content)
	if err := uc.messageRepo.Create(ctx, userMessage); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}
	return userMessage, nil
}

// getConversationHistory retrieves conversation history in LLM format
func (uc *ChatUseCase) getConversationHistory(ctx context.Context, session *entity.Session) ([]ports.Message, error) {
	messages, err := uc.messageRepo.FindBySessionID(ctx, string(session.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	llmMessages := make([]ports.Message, 0, len(messages))
	for _, msg := range messages {
		llmMessages = append(llmMessages, ports.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}
	return llmMessages, nil
}

// saveAssistantMessage saves assistant message to repository
func (uc *ChatUseCase) saveAssistantMessage(ctx context.Context, session *entity.Session, content string) (*entity.Message, error) {
	assistantMessage := entity.NewAssistantMessage(string(session.ID), content)
	if err := uc.messageRepo.Create(ctx, assistantMessage); err != nil {
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}
	return assistantMessage, nil
}

// updateSession updates session timestamp
func (uc *ChatUseCase) updateSession(ctx context.Context, session *entity.Session) error {
	session.UpdateTimestamp()
	return uc.sessionRepo.Update(ctx, session)
}

// buildSendMessageResponse builds response with updated conversation history
func (uc *ChatUseCase) buildSendMessageResponse(ctx context.Context, session *entity.Session, assistantMessage *entity.Message) (*dto.SendMessageResponse, error) {
	// Get updated messages for response
	updatedMessages, err := uc.messageRepo.FindBySessionID(ctx, string(session.ID))
	if err != nil {
		uc.logger.Error("failed to get updated messages", "error", err)
	}

	// Convert messages to DTOs
	messageDTOs := make([]*dto.MessageDTO, 0, len(updatedMessages))
	for _, msg := range updatedMessages {
		messageDTOs = append(messageDTOs, dto.MessageDTOFromEntity(msg))
	}

	return &dto.SendMessageResponse{
		Success:  true,
		Message:  dto.MessageDTOFromEntity(assistantMessage),
		Messages: messageDTOs,
	}, nil
}
