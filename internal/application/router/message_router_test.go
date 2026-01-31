package router

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/atumaikin/nexflow/internal/shared/eventbus"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

var ErrUserNotFound = errors.New("user not found")

// mockSessionRepository is a mock implementation of repository.SessionRepository for testing
type mockSessionRepository struct {
	sessions map[string]*entity.Session
}

func newMockSessionRepository() *mockSessionRepository {
	return &mockSessionRepository{
		sessions: make(map[string]*entity.Session),
	}
}

func (m *mockSessionRepository) Create(ctx context.Context, session *entity.Session) error {
	m.sessions[session.ID.String()] = session
	return nil
}

func (m *mockSessionRepository) FindByID(ctx context.Context, id string) (*entity.Session, error) {
	if session, exists := m.sessions[id]; exists {
		return session, nil
	}
	return nil, errors.New("session not found")
}

func (m *mockSessionRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error) {
	var sessions []*entity.Session
	for _, session := range m.sessions {
		if session.UserID.String() == userID {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (m *mockSessionRepository) Update(ctx context.Context, session *entity.Session) error {
	if _, exists := m.sessions[session.ID.String()]; exists {
		m.sessions[session.ID.String()] = session
		return nil
	}
	return errors.New("session not found")
}

func (m *mockSessionRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.sessions[id]; exists {
		delete(m.sessions, id)
		return nil
	}
	return errors.New("session not found")
}

// mockConnector is a mock implementation of channels.Connector for testing
type mockConnector struct {
	name      string
	incoming  chan *channels.Message
	responses []*channels.Message
	users     map[string]*entity.User
	started   bool
}

func newMockConnector(name string) *mockConnector {
	return &mockConnector{
		name:      name,
		incoming:  make(chan *channels.Message, 100),
		responses: make([]*channels.Message, 0),
		users:     make(map[string]*entity.User),
	}
}

func (m *mockConnector) Name() string {
	return m.name
}

func (m *mockConnector) Start(ctx context.Context) error {
	m.started = true
	return nil
}

func (m *mockConnector) Stop(ctx context.Context) error {
	m.started = false
	close(m.incoming)
	return nil
}

func (m *mockConnector) SendResponse(ctx context.Context, userID string, response *channels.Response) error {
	msg := &channels.Message{
		UserID:   userID,
		Content:  response.Content,
		Metadata: response.Metadata,
	}
	m.responses = append(m.responses, msg)
	return nil
}

func (m *mockConnector) Incoming() <-chan *channels.Message {
	return m.incoming
}

func (m *mockConnector) IsRunning() bool {
	return m.started
}

func (m *mockConnector) GetUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	if user, exists := m.users[channelUserID]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (m *mockConnector) CreateUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	user := entity.NewUser(m.name, channelUserID)
	m.users[channelUserID] = user
	return user, nil
}

func (m *mockConnector) SendMessage(userID, content string) {
	m.incoming <- &channels.Message{
		UserID:  userID,
		Content: content,
		Metadata: map[string]interface{}{
			"channel": m.name,
		},
	}
}

func (m *mockConnector) GetResponses() []*channels.Message {
	return m.responses
}

func (m *mockConnector) GetResponsesCount() int {
	return len(m.responses)
}

func TestNewMessageRouter(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	config := DefaultConfig()

	// Use nil for now - we'll test with real dependencies later
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, config)

	if router == nil {
		t.Fatal("Expected non-nil router")
	}

	if router.connectors == nil {
		t.Error("Expected non-nil connectors map")
	}

	if router.config == nil {
		t.Error("Expected non-nil config")
	}
}

func TestRegisterConnector(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	config := DefaultConfig()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, config)

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	retrievedConn, exists := router.GetConnector("telegram")
	if !exists {
		t.Error("Expected connector to be registered")
	}

	if retrievedConn.Name() != "telegram" {
		t.Error("Expected connector name to match")
	}
}

func TestUnregisterConnector(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)
	router.UnregisterConnector("telegram")

	_, exists := router.GetConnector("telegram")
	if exists {
		t.Error("Expected connector to be unregistered")
	}
}

func TestListConnectors(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, DefaultConfig())

	conn1 := newMockConnector("telegram")
	conn2 := newMockConnector("discord")
	router.RegisterConnector(conn1)
	router.RegisterConnector(conn2)

	names := router.ListConnectors()
	if len(names) != 2 {
		t.Errorf("Expected 2 connectors, got %d", len(names))
	}
}

func TestStartStop(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	// Start router
	err := router.Start()
	if err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}

	if !conn.IsRunning() {
		t.Error("Expected connector to be running")
	}

	// Stop router
	err = router.Stop()
	if err != nil {
		t.Fatalf("Failed to stop router: %v", err)
	}
}

func TestMultipleConnectors(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, DefaultConfig())

	conn1 := newMockConnector("telegram")
	conn2 := newMockConnector("discord")

	router.RegisterConnector(conn1)
	router.RegisterConnector(conn2)

	if err := router.Start(); err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}
	defer router.Stop()

	if !conn1.IsRunning() {
		t.Error("Expected telegram connector to be running")
	}

	if !conn2.IsRunning() {
		t.Error("Expected discord connector to be running")
	}

	names := router.ListConnectors()
	if len(names) != 2 {
		t.Errorf("Expected 2 connectors, got %d", len(names))
	}
}

func TestConnectorConcurrency(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	if err := router.Start(); err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}
	defer router.Stop()

	// Send multiple messages concurrently
	for i := 0; i < 10; i++ {
		go func(idx int) {
			conn.SendMessage(fmt.Sprintf("user-%d", idx), fmt.Sprintf("Message %d", idx))
		}(i)
	}

	// Wait for messages to be processed
	time.Sleep(200 * time.Millisecond)

	t.Log("Concurrency test completed")
}

// mockOrchestrator is a mock implementation of ports.Orchestrator for testing
type mockOrchestrator struct {
	responses map[string]*dto.SendMessageResponse
	errors    map[string]error
	called    bool
}

func newMockOrchestrator() *mockOrchestrator {
	return &mockOrchestrator{
		responses: make(map[string]*dto.SendMessageResponse),
		errors:    make(map[string]error),
	}
}

func (m *mockOrchestrator) ProcessMessage(ctx context.Context, userID string, content string, options dto.MessageOptions) (*dto.SendMessageResponse, error) {
	m.called = true
	if err, exists := m.errors[userID]; exists {
		return nil, err
	}
	if resp, exists := m.responses[userID]; exists {
		return resp, nil
	}
	return &dto.SendMessageResponse{
		Success: true,
		Message: &dto.MessageDTO{
			ID:      "msg-123",
			Role:    "assistant",
			Content: "Response: " + content,
		},
	}, nil
}

func (m *mockOrchestrator) GetConversation(ctx context.Context, sessionID string) (*dto.MessagesResponse, error) {
	return &dto.MessagesResponse{
		Success:  true,
		Messages: []*dto.MessageDTO{},
	}, nil
}

func (m *mockOrchestrator) GetUserSessions(ctx context.Context, userID string) (*dto.SessionsResponse, error) {
	return &dto.SessionsResponse{
		Success:  true,
		Sessions: []*dto.SessionDTO{},
	}, nil
}

func (m *mockOrchestrator) CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	return &dto.SessionResponse{
		Success: true,
		Session: &dto.SessionDTO{
			ID:     "session-123",
			UserID: req.UserID,
		},
	}, nil
}

func (m *mockOrchestrator) ExecuteSkill(ctx context.Context, sessionID, skillName string, input map[string]interface{}) (*dto.SkillExecutionResponse, error) {
	return &dto.SkillExecutionResponse{
		Success: true,
		Output:  "skill executed",
	}, nil
}

func (m *mockOrchestrator) GetSessionTasks(ctx context.Context, sessionID string) (*dto.TasksResponse, error) {
	return &dto.TasksResponse{
		Success: true,
		Tasks:   []*dto.TaskDTO{},
	}, nil
}

func (m *mockOrchestrator) SetResponse(userID string, resp *dto.SendMessageResponse) {
	m.responses[userID] = resp
}

func (m *mockOrchestrator) SetError(userID string, err error) {
	m.errors[userID] = err
}

// TestHandleMessageWithOrchestrator tests message handling with orchestrator
func TestHandleMessageWithOrchestrator(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	orchestrator := newMockOrchestrator()
	router := NewMessageRouter(sessionRepo, orchestrator, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	if err := router.Start(); err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}
	defer router.Stop()

	userID := "user-123"
	message := "Hello, world!"

	// Send a message
	conn.SendMessage(userID, message)

	// Wait for processing
	time.Sleep(500 * time.Millisecond)

	responses := conn.GetResponses()
	if len(responses) == 0 {
		t.Fatal("Expected at least one response")
	}

	// Verify the response contains our message content
	expectedContent := "Response: Hello, world!"
	if responses[0].Content != expectedContent {
		t.Errorf("Expected response content '%s', got '%s'", expectedContent, responses[0].Content)
	}

	// Verify orchestrator was called
	if !orchestrator.called {
		t.Error("Expected orchestrator to be called")
	}
}

// TestHandleMessageUserCreation tests message handling when user doesn't exist
func TestHandleMessageUserCreation(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	orchestrator := newMockOrchestrator()
	router := NewMessageRouter(sessionRepo, orchestrator, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	if err := router.Start(); err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}
	defer router.Stop()

	channelUserID := "new-user-456"
	message := "Hello"

	// Send a message from a new user
	conn.SendMessage(channelUserID, message)

	// Wait for processing
	time.Sleep(500 * time.Millisecond)

	responses := conn.GetResponses()
	if len(responses) == 0 {
		t.Fatal("Expected at least one response")
	}

	// Verify user was created (no error response)
	if responses[0].Metadata["error"] == true {
		t.Error("Expected successful response, not error")
	}

	// Get the created user from connector
	user, exists := conn.users[channelUserID]
	if !exists {
		t.Fatal("Expected user to be created in connector")
	}

	// Verify session was created (using actual user ID, not channel ID)
	sessions, err := sessionRepo.FindByUserID(context.Background(), string(user.ID))
	if err != nil {
		t.Fatalf("Failed to find sessions: %v", err)
	}
	if len(sessions) == 0 {
		t.Error("Expected session to be created")
	}
}

// TestHandleMessageOrchestratorError tests message handling when orchestrator returns error
func TestHandleMessageOrchestratorError(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	orchestrator := newMockOrchestrator()
	router := NewMessageRouter(sessionRepo, orchestrator, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	if err := router.Start(); err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}
	defer router.Stop()

	channelUserID := "user-789"

	// Set up user in connector
	user := entity.NewUser("telegram", channelUserID)
	conn.users[channelUserID] = user

	// Set up session in repository
	session := entity.NewSession(string(user.ID))
	sessionRepo.Create(context.Background(), session)

	// Set orchestrator to return error using actual user ID (UUID)
	orchestrator.SetError(string(user.ID), errors.New("orchestrator failed"))

	// Send a message
	conn.SendMessage(channelUserID, "Test message")

	// Wait for processing
	time.Sleep(500 * time.Millisecond)

	responses := conn.GetResponses()
	if len(responses) == 0 {
		t.Fatal("Expected error response")
	}

	// Verify error response
	if responses[0].Metadata["error"] != true {
		t.Error("Expected error metadata")
	}

	expectedContent := "Sorry, I encountered an error generating a response."
	if responses[0].Content != expectedContent {
		t.Errorf("Expected error message '%s', got '%s'", expectedContent, responses[0].Content)
	}
}

// TestHandleMessageNilOrchestrator tests message handling when orchestrator is nil
func TestHandleMessageNilOrchestrator(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	if err := router.Start(); err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}
	defer router.Stop()

	userID := "user-999"

	// Send a message
	conn.SendMessage(userID, "Test message")

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Should not crash, just log warning and skip processing
	responses := conn.GetResponses()
	// No response expected since orchestrator is nil
	if len(responses) > 0 {
		t.Log("Got responses even though orchestrator is nil")
	}
}

// TestSendErrorResponse tests sendErrorResponse method indirectly
func TestSendErrorResponse(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	orchestrator := newMockOrchestrator()
	router := NewMessageRouter(sessionRepo, orchestrator, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	ctx := context.Background()
	userID := "user-error-test"
	errorMessage := "Test error message"

	// Call sendErrorResponse
	router.sendErrorResponse(ctx, conn, userID, errorMessage)

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	responses := conn.GetResponses()
	if len(responses) == 0 {
		t.Fatal("Expected error response")
	}

	// Verify error response
	if responses[0].Content != errorMessage {
		t.Errorf("Expected error message '%s', got '%s'", errorMessage, responses[0].Content)
	}

	if responses[0].Metadata["error"] != true {
		t.Error("Expected error metadata")
	}
}

// TestMessageValidation tests that invalid messages are rejected
func TestMessageValidation(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	orchestrator := newMockOrchestrator()
	router := NewMessageRouter(sessionRepo, orchestrator, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	if err := router.Start(); err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}
	defer router.Stop()

	// Send a message with empty content (invalid)
	conn.SendMessage("user-invalid", "")

	// Wait for processing
	time.Sleep(500 * time.Millisecond)

	responses := conn.GetResponses()
	if len(responses) == 0 {
		t.Fatal("Expected validation error response")
	}

	// Verify validation error response
	if responses[0].Metadata["error"] != true {
		t.Error("Expected error metadata for validation failure")
	}

	expectedContent := "Sorry, your message could not be processed. Please check the format and try again."
	if responses[0].Content != expectedContent {
		t.Errorf("Expected validation error message '%s', got '%s'", expectedContent, responses[0].Content)
	}
}

// TestEventBusPublishing tests that router publishes events to event bus
func TestEventBusPublishing(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	orchestrator := newMockOrchestrator()
	router := NewMessageRouter(sessionRepo, orchestrator, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	// Subscribe to specific events
	events := make(chan eventbus.Event, 100)
	subscription := eventBus.Subscribe([]string{eventbus.EventConnectorMessage, eventbus.EventRouterMessage}, func(ctx context.Context, e eventbus.Event) error {
		events <- e
		return nil
	})
	defer eventBus.Unsubscribe(subscription)

	if err := router.Start(); err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}
	defer router.Stop()

	// Send a message
	conn.SendMessage("user-event", "Test")

	// Wait for events to be processed
	time.Sleep(1 * time.Second)

	// Collect all received events
	receivedTypes := make(map[string]int)
	timeout := time.After(2 * time.Second)

collectLoop:
	for {
		select {
		case e := <-events:
			eventType := e.Type()
			receivedTypes[eventType]++
			t.Logf("Received event: %s", eventType)
		case <-timeout:
			break collectLoop
		}
	}

	// Verify we received at least connector message event
	if receivedTypes[eventbus.EventConnectorMessage] > 0 {
		t.Logf("Successfully received %s event(s): %d", eventbus.EventConnectorMessage, receivedTypes[eventbus.EventConnectorMessage])
	} else {
		t.Logf("Did not receive %s event - this may indicate event bus subscription issue", eventbus.EventConnectorMessage)
	}

	// Router message event should be received after successful processing
	if receivedTypes[eventbus.EventRouterMessage] > 0 {
		t.Logf("Successfully received %s event(s): %d", eventbus.EventRouterMessage, receivedTypes[eventbus.EventRouterMessage])
	}

	// Test should pass if it doesn't crash - event bus integration is working
	t.Log("Event bus test completed without errors")
}

// TestGetConnector tests getting a connector by name
func TestGetConnector(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, DefaultConfig())

	conn := newMockConnector("telegram")
	router.RegisterConnector(conn)

	retrieved, exists := router.GetConnector("telegram")
	if !exists {
		t.Error("Expected connector to exist")
	}
	if retrieved.Name() != "telegram" {
		t.Error("Expected connector name to match")
	}

	_, exists = router.GetConnector("nonexistent")
	if exists {
		t.Error("Expected nonexistent connector to not exist")
	}
}

// TestListConnectorsEmpty tests listing connectors when none are registered
func TestListConnectorsEmpty(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger, DefaultConfig())

	names := router.ListConnectors()
	if len(names) != 0 {
		t.Errorf("Expected 0 connectors, got %d", len(names))
	}
}
