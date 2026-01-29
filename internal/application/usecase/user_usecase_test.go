package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/shared/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByChannel(ctx context.Context, channel, channelID string) (*entity.User, error) {
	args := m.Called(ctx, channel, channelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*entity.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockLogger is a mock implementation of Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	m.Called(ctx, msg, args)
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	m.Called(ctx, msg, args)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	m.Called(ctx, msg, args)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	m.Called(ctx, msg, args)
}

func (m *MockLogger) With(args ...interface{}) logging.Logger {
	return &MockLogger{}
}

func (m *MockLogger) WithContext(ctx context.Context) logging.Logger {
	return &MockLogger{}
}

func TestUserUseCase_CreateUser_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	req := dto.CreateUserRequest{
		Channel:   "telegram",
		ChannelID: "user123",
	}

	mockRepo.On("FindByChannel", ctx, req.Channel, req.ChannelID).Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
	mockLogger.On("Info", "user created", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// Act
	resp, err := uc.CreateUser(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.User)
	assert.Equal(t, req.Channel, resp.User.Channel)
	assert.Equal(t, req.ChannelID, resp.User.ChannelID)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_CreateUser_AlreadyExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	existingUser := entity.NewUser("telegram", "user123")
	req := dto.CreateUserRequest{
		Channel:   "telegram",
		ChannelID: "user123",
	}

	mockRepo.On("FindByChannel", ctx, req.Channel, req.ChannelID).Return(existingUser, nil)

	// Act
	resp, err := uc.CreateUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.False(t, resp.Success)
	assert.NotEmpty(t, resp.Error)
	assert.Contains(t, resp.Error, "user already exists")
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_CreateUser_RepoError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	req := dto.CreateUserRequest{
		Channel:   "telegram",
		ChannelID: "user123",
	}

	mockRepo.On("FindByChannel", ctx, req.Channel, req.ChannelID).Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))

	// Act
	resp, err := uc.CreateUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.Error, "failed to create user")
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetUserByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	user := entity.NewUser("telegram", "user123")
	userID := string(user.ID)

	mockRepo.On("FindByID", ctx, userID).Return(user, nil)

	// Act
	resp, err := uc.GetUserByID(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.User)
	assert.Equal(t, userID, resp.User.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetUserByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	mockRepo.On("FindByID", ctx, "nonexistent").Return(nil, errors.New("user not found"))

	// Act
	resp, err := uc.GetUserByID(ctx, "nonexistent")

	// Assert
	require.Error(t, err)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.Error, "failed to find user")
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetUserByChannel_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	user := entity.NewUser("telegram", "user123")

	mockRepo.On("FindByChannel", ctx, "telegram", "user123").Return(user, nil)

	// Act
	resp, err := uc.GetUserByChannel(ctx, "telegram", "user123")

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.User)
	assert.Equal(t, "telegram", resp.User.Channel)
	assert.Equal(t, "user123", resp.User.ChannelID)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_ListUsers_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	users := []*entity.User{
		entity.NewUser("telegram", "user1"),
		entity.NewUser("discord", "user2"),
		entity.NewUser("web", "user3"),
	}

	mockRepo.On("List", ctx).Return(users, nil)

	// Act
	resp, err := uc.ListUsers(ctx)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, resp.Users, 3)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_ListUsers_Empty(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	mockRepo.On("List", ctx).Return([]*entity.User{}, nil)

	// Act
	resp, err := uc.ListUsers(ctx)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Empty(t, resp.Users)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_DeleteUser_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	user := entity.NewUser("telegram", "user123")
	userID := string(user.ID)

	mockRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockRepo.On("Delete", ctx, userID).Return(nil)
	mockLogger.On("Info", "user deleted", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// Act
	resp, err := uc.DeleteUser(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.User)
	assert.Equal(t, userID, resp.User.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_DeleteUser_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	mockRepo.On("FindByID", ctx, "nonexistent").Return(nil, errors.New("user not found"))

	// Act
	resp, err := uc.DeleteUser(ctx, "nonexistent")

	// Assert
	require.Error(t, err)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.Error, "user not found")
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetOrCreateUser_Existing(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	existingUser := entity.NewUser("telegram", "user123")

	mockRepo.On("FindByChannel", ctx, "telegram", "user123").Return(existingUser, nil)

	// Act
	resp, err := uc.GetOrCreateUser(ctx, "telegram", "user123")

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.User)
	assert.Equal(t, string(existingUser.ID), resp.User.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetOrCreateUser_NewUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	uc := NewUserUseCase(mockRepo, mockLogger)

	mockRepo.On("FindByChannel", ctx, "telegram", "user123").Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
	mockLogger.On("Info", "user created", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// Act
	resp, err := uc.GetOrCreateUser(ctx, "telegram", "user123")

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.User)
	assert.Equal(t, "telegram", resp.User.Channel)
	assert.Equal(t, "user123", resp.User.ChannelID)
	mockRepo.AssertExpectations(t)
}
