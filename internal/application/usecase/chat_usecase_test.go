package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSessionRepository is a mock implementation of SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) FindByID(ctx context.Context, id string) (*entity.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) Update(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockMessageRepository is a mock implementation of MessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(ctx context.Context, message *entity.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) FindByID(ctx context.Context, id string) (*entity.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Message), args.Error(1)
}

func (m *MockMessageRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Message, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Message), args.Error(1)
}

func (m *MockMessageRepository) Update(ctx context.Context, message *entity.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMessageRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

// MockTaskRepository is a mock implementation of TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id string) (*entity.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Task, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockLLMProvider is a mock implementation of LLMProvider
type MockLLMProvider struct {
	mock.Mock
}

func (m *MockLLMProvider) Generate(ctx context.Context, req ports.CompletionRequest) (*ports.CompletionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.CompletionResponse), args.Error(1)
}

func (m *MockLLMProvider) GenerateWithTools(ctx context.Context, req ports.CompletionRequest, tools []ports.ToolDefinition) (*ports.CompletionResponse, error) {
	args := m.Called(ctx, req, tools)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.CompletionResponse), args.Error(1)
}

func (m *MockLLMProvider) Stream(ctx context.Context, req ports.CompletionRequest) (<-chan string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(<-chan string), args.Error(1)
}

func (m *MockLLMProvider) EstimateCost(req ports.CompletionRequest) (float64, error) {
	args := m.Called(req)
	return args.Get(0).(float64), args.Error(1)
}

func TestChatUseCase_SendMessage_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockTaskRepo := new(MockTaskRepository)
	mockLLMProvider := new(MockLLMProvider)
	mockLogger := new(MockLogger)

	uc := NewChatUseCase(mockUserRepo, mockSessionRepo, mockMessageRepo, mockTaskRepo, mockLLMProvider, mockLogger)

	user := entity.NewUser("web", "user123")
	session := entity.NewSession(user.ID)
	userMsg := entity.NewUserMessage(session.ID, "Hello")

	req := dto.SendMessageRequest{
		UserID: "user123",
		Message: dto.ChatMessage{
			Role:    "user",
			Content: "Hello",
		},
		Options: dto.MessageOptions{
			Model:     "gpt-4",
			MaxTokens: 1000,
		},
	}

	llmResp := &ports.CompletionResponse{
		Message: ports.Message{
			Role:    "assistant",
			Content: "Hi there!",
		},
		Tokens: ports.Tokens{
			InputTokens:  10,
			OutputTokens: 5,
			TotalTokens:  15,
		},
	}

	assistantMsg := entity.NewAssistantMessage(session.ID, "Hi there!")

	mockUserRepo.On("FindByChannel", ctx, "web", "user123").Return(nil, errors.New("not found"))
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)
	mockMessageRepo.On("Create", ctx, mock.AnythingOfType("*entity.Message")).Return(nil)
	mockMessageRepo.On("FindBySessionID", ctx, mock.Anything).Return([]*entity.Message{userMsg}, nil)
	mockLLMProvider.On("Generate", ctx, mock.AnythingOfType("ports.CompletionRequest")).Return(llmResp, nil)
	mockMessageRepo.On("Create", ctx, mock.AnythingOfType("*entity.Message")).Return(nil)
	mockSessionRepo.On("Update", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)
	mockMessageRepo.On("FindBySessionID", ctx, mock.Anything).Return([]*entity.Message{userMsg, assistantMsg}, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return().Maybe()

	// Act
	resp, err := uc.SendMessage(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Message)
	assert.Equal(t, "assistant", resp.Message.Role)
	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
	mockLLMProvider.AssertExpectations(t)
}

func TestChatUseCase_GetConversation_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockTaskRepo := new(MockTaskRepository)
	mockLLMProvider := new(MockLLMProvider)
	mockLogger := new(MockLogger)

	uc := NewChatUseCase(mockUserRepo, mockSessionRepo, mockMessageRepo, mockTaskRepo, mockLLMProvider, mockLogger)

	sessionID := "session-1"
	messages := []*entity.Message{
		entity.NewUserMessage(sessionID, "Hello"),
		entity.NewAssistantMessage(sessionID, "Hi there!"),
	}

	mockMessageRepo.On("FindBySessionID", ctx, sessionID).Return(messages, nil)

	// Act
	resp, err := uc.GetConversation(ctx, sessionID)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Messages, 2)
	mockMessageRepo.AssertExpectations(t)
}

func TestChatUseCase_GetUserSessions_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockTaskRepo := new(MockTaskRepository)
	mockLLMProvider := new(MockLLMProvider)
	mockLogger := new(MockLogger)

	uc := NewChatUseCase(mockUserRepo, mockSessionRepo, mockMessageRepo, mockTaskRepo, mockLLMProvider, mockLogger)

	userID := "user-1"
	sessions := []*entity.Session{
		entity.NewSession(userID),
		entity.NewSession(userID),
	}

	mockSessionRepo.On("FindByUserID", ctx, userID).Return(sessions, nil)

	// Act
	resp, err := uc.GetUserSessions(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Sessions, 2)
	mockSessionRepo.AssertExpectations(t)
}

func TestChatUseCase_CreateSession_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockTaskRepo := new(MockTaskRepository)
	mockLLMProvider := new(MockLLMProvider)
	mockLogger := new(MockLogger)

	uc := NewChatUseCase(mockUserRepo, mockSessionRepo, mockMessageRepo, mockTaskRepo, mockLLMProvider, mockLogger)

	req := dto.CreateSessionRequest{
		UserID: "user-1",
	}

	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)

	// Act
	resp, err := uc.CreateSession(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Session)
	assert.Equal(t, "user-1", resp.Session.UserID)
	mockSessionRepo.AssertExpectations(t)
}

func TestChatUseCase_ExecuteSkill_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockTaskRepo := new(MockTaskRepository)
	mockLLMProvider := new(MockLLMProvider)
	mockLogger := new(MockLogger)

	uc := NewChatUseCase(mockUserRepo, mockSessionRepo, mockMessageRepo, mockTaskRepo, mockLLMProvider, mockLogger)

	sessionID := "session-1"
	skillName := "my-skill"
	input := map[string]interface{}{"param": "value"}

	mockTaskRepo.On("Create", ctx, mock.AnythingOfType("*entity.Task")).Return(nil)
	mockTaskRepo.On("Update", ctx, mock.AnythingOfType("*entity.Task")).Return(nil)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return().Maybe()

	// Act
	resp, err := uc.ExecuteSkill(ctx, sessionID, skillName, input)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.Output, "skill my-skill executed successfully")
	mockTaskRepo.AssertExpectations(t)
}

func TestChatUseCase_GetSessionTasks_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockTaskRepo := new(MockTaskRepository)
	mockLLMProvider := new(MockLLMProvider)
	mockLogger := new(MockLogger)

	uc := NewChatUseCase(mockUserRepo, mockSessionRepo, mockMessageRepo, mockTaskRepo, mockLLMProvider, mockLogger)

	sessionID := "session-1"
	task1 := entity.NewTask(sessionID, "skill1", "{}")
	task2 := entity.NewTask(sessionID, "skill2", "{}")
	task2.SetCompleted(`{"result": "ok"}`)

	tasks := []*entity.Task{task1, task2}

	mockTaskRepo.On("FindBySessionID", ctx, sessionID).Return(tasks, nil)

	// Act
	resp, err := uc.GetSessionTasks(ctx, sessionID)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Tasks, 2)
	mockTaskRepo.AssertExpectations(t)
}
