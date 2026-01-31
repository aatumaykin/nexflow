package service

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
)

// SessionAccessService provides access control for sessions.
// It checks whether a user can access a specific session based on ownership and permissions.
type SessionAccessService struct {
	sessionRepo repository.SessionRepository
}

// NewSessionAccessService creates a new SessionAccessService.
//
// Parameters:
//   - sessionRepo: SessionRepository for retrieving session data
//
// Returns:
//   - *SessionAccessService: Initialized session access service
func NewSessionAccessService(sessionRepo repository.SessionRepository) *SessionAccessService {
	return &SessionAccessService{
		sessionRepo: sessionRepo,
	}
}

// CanAccessSession checks if a user can access the specified session.
// Users can only access their own sessions.
//
// Parameters:
//   - ctx: Context for the operation
//   - userID: ID of the user requesting access
//   - sessionID: ID of the session to access
//
// Returns:
//   - bool: True if the user can access the session, false otherwise
func (s *SessionAccessService) CanAccessSession(ctx context.Context, userID valueobject.UserID, sessionID valueobject.SessionID) bool {
	session, err := s.sessionRepo.FindByID(ctx, string(sessionID))
	if err != nil || session == nil {
		// Session not found, error occurred, or session is nil - deny access
		return false
	}

	// Check if session belongs to the user
	return session.IsOwnedBy(userID)
}
