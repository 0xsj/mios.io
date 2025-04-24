package token

import (
	"fmt"
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
}

type SessionManager struct{}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func (m *SessionManager) CreateSession(userID, userAgent, ip string, duration time.Duration) (*Session, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	now := time.Now()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(duration),
		UserAgent: userAgent,
		IP:        ip,
	}

	// save session here

	return session, nil
}

func (m *SessionManager) ValidateSession(sessionID string) (*Session, error) {
	// fetch session
	// check if session expired
	return nil, fmt.Errorf("not implemented: requires session repository")
}

func (m *SessionManager) RevokeSession(sessionID string) error {
	// delete session
	return fmt.Errorf("not implemented: requires session repository")
}

func (m *SessionManager) ExtendedSession(sessionID string, duration time.Duration) error {
	// update session
	return fmt.Errorf("not implemented: requires session repository")
}
