package router

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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

	// Use nil for now - we'll test with real dependencies later
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger)

	if router == nil {
		t.Fatal("Expected non-nil router")
	}

	if router.connectors == nil {
		t.Error("Expected non-nil connectors map")
	}
}

func TestRegisterConnector(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	sessionRepo := newMockSessionRepository()
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger)

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
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger)

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
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger)

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
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger)

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
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger)

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
	router := NewMessageRouter(sessionRepo, nil, eventBus, logger)

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
