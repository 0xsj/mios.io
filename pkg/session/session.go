package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/redis"
	"github.com/0xsj/gin-sqlc/pkg/token"
)

const (
	SessionPrefix     = "session:"
	DefaultExpiration = 24 * time.Hour
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Store struct {
	client *redis.Client
	logger log.Logger
}

func NewStore(client *redis.Client, logger log.Logger) *Store {
	return &Store{
		client: client,
		logger: logger,
	}
}

func (s *Store) Create(ctx context.Context, userID, userAgent, ip string, expiration time.Duration) (*Session, error) {
	if expiration == 0 {
		expiration = DefaultExpiration
	}

	sessionID, err := token.GenerateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	now := time.Now()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		UserAgent: userAgent,
		IP:        ip,
		CreatedAt: now,
		ExpiresAt: now.Add(expiration),
	}

	data, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session: %w", err)
	}

	key := SessionPrefix + sessionID
	err = s.client.Set(ctx, key, data, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	userKey := fmt.Sprintf("user:%s:sessions", userID)
	err = s.client.Set(ctx, userKey+":"+sessionID, sessionID, expiration)
	if err != nil {
		s.logger.Warnf("Failed to store user-session mapping: %v", err)
	}

	s.logger.Infof("Created session %s for user %s", sessionID, userID)
	return session, nil
}

func (s *Store) Get(ctx context.Context, sessionID string) (*Session, error) {
	key := SessionPrefix + sessionID
	data, err := s.client.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session Session
	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		// FIX: Handle the error from Delete
		if err := s.Delete(ctx, sessionID); err != nil {
			s.logger.Warnf("Failed to delete expired session: %v", err)
			// Continue despite the error since we're just cleaning up
		}
		return nil, nil
	}

	return &session, nil
}

func (s *Store) Delete(ctx context.Context, sessionID string) error {
	key := SessionPrefix + sessionID

	// Get session first to get user ID
	session, err := s.Get(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session for deletion: %w", err)
	}
	if session != nil {
		// Remove user-session mapping
		userKey := fmt.Sprintf("user:%s:sessions:%s", session.UserID, sessionID)
		if err := s.client.Delete(ctx, userKey); err != nil {
			s.logger.Warnf("Failed to delete user-session mapping: %v", err)
			// Non-critical error, continue
		}
	}

	// Remove the session
	err = s.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	s.logger.Infof("Deleted session %s", sessionID)
	return nil
}

func (s *Store) DeleteByUserID(ctx context.Context, userID string) error {
	userPattern := fmt.Sprintf("user:%s:sessions:*", userID)

	// Find all sessions for this user
	keys, err := s.client.Keys(ctx, userPattern)
	if err != nil {
		return fmt.Errorf("failed to find user sessions: %w", err)
	}

	// Delete each session
	for _, key := range keys {
		// Extract session ID from key
		sessionID := key[len(fmt.Sprintf("user:%s:sessions:", userID)):]
		err := s.Delete(ctx, sessionID)
		if err != nil {
			s.logger.Warnf("Failed to delete session %s: %v", sessionID, err)
			// Continue with other sessions
		}
	}

	s.logger.Infof("Deleted all sessions for user %s", userID)
	return nil
}

// Refresh extends a session's expiration
func (s *Store) Refresh(ctx context.Context, sessionID string, expiration time.Duration) error {
	if expiration == 0 {
		expiration = DefaultExpiration
	}

	key := SessionPrefix + sessionID

	// Get the current session
	session, err := s.Get(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session for refresh: %w", err)
	}
	if session == nil {
		return fmt.Errorf("session not found")
	}

	// Update expiration
	session.ExpiresAt = time.Now().Add(expiration)

	// Serialize and store the updated session
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	err = s.client.Set(ctx, key, data, expiration)
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	// Also refresh the user-session mapping
	userKey := fmt.Sprintf("user:%s:sessions:%s", session.UserID, sessionID)
	err = s.client.Expire(ctx, userKey, expiration)
	if err != nil {
		s.logger.Warnf("Failed to refresh user-session mapping: %v", err)
		// Non-critical error, continue
	}

	s.logger.Infof("Refreshed session %s for user %s", sessionID, session.UserID)
	return nil
}
