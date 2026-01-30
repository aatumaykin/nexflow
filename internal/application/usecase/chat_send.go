package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/ports"
)

// SendMessage processes a user message and returns AI response
//
// Parameters:
//   - ctx: Context for the operation
//   - req: SendMessageRequest containing message and LLM options
//
// Returns:
//   - *dto.SendMessageResponse: Response containing AI message and conversation history
//   - error: Error if operation failed
func (uc *ChatUseCase) SendMessage(ctx context.Context, req dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	user, err := uc.findOrCreateUser(ctx, req.UserID)
	if err != nil {
		return handleSendError(err, "failed to get user")
	}

	session, err := uc.createSession(ctx, user)
	if err != nil {
		return handleSendError(err, "failed to create session")
	}

	_, err = uc.saveUserMessage(ctx, session, req.Message.Content)
	if err != nil {
		return handleSendError(err, "failed to save user message")
	}

	llmMessages, err := uc.getConversationHistory(ctx, session)
	if err != nil {
		return handleSendError(err, "failed to get conversation history")
	}

	llmResp, err := uc.callLLM(ctx, llmMessages, req.Options)
	if err != nil {
		return handleSendError(err, "failed to generate response")
	}

	// Save assistant message
	assistantMessage, err := uc.saveAssistantMessage(ctx, session, llmResp.Message.Content)
	if err != nil {
		uc.logger.Error("failed to save assistant message", "error", err)
	}

	// Update session
	if err := uc.updateSession(ctx, session); err != nil {
		uc.logger.Error("failed to update session", "error", err)
	}

	// Build response
	return uc.buildSendMessageResponse(ctx, session, assistantMessage)
}

// callLLM calls LLM provider with conversation history
func (uc *ChatUseCase) callLLM(ctx context.Context, messages []ports.Message, options dto.MessageOptions) (*ports.CompletionResponse, error) {
	llmReq := ports.CompletionRequest{
		Messages:  messages,
		Model:     options.Model,
		MaxTokens: options.MaxTokens,
	}
	return uc.llmProvider.Generate(ctx, llmReq)
}
