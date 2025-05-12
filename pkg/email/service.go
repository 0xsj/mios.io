// pkg/email/service.go
package email

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/google/uuid"
)

// Config holds email service configuration
type Config struct {
	DefaultFrom     string
	DefaultReplyTo  string
	TemplateDir     string
	MaxRetries      int
	RetryInterval   time.Duration
	Providers       []Provider
	DefaultProvider string
	Logger          log.Logger
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		MaxRetries:    3,
		RetryInterval: time.Second * 10,
	}
}

// emailService implements the Service interface
type emailService struct {
	config      Config
	logger      log.Logger
	providers   map[string]Provider
	defaultProvider Provider
	templates   map[string]string
	mu          sync.RWMutex
	queue       map[string]*queuedMessage
}

type queuedMessage struct {
	id         string
	message    *Message
	provider   Provider
	status     MessageStatus
	retries    int
	lastAttempt time.Time
	errorMsg   string
}

// NewEmailService creates a new email service
func NewEmailService(config Config) (Service, error) {
	if len(config.Providers) == 0 {
		return nil, errors.New("at least one email provider is required")
	}
	
	s := &emailService{
		config:    config,
		logger:    config.Logger,
		providers: make(map[string]Provider),
		templates: make(map[string]string),
		queue:     make(map[string]*queuedMessage),
	}
	
	// Register providers
	for _, provider := range config.Providers {
		s.providers[provider.Name()] = provider
		if provider.Name() == config.DefaultProvider {
			s.defaultProvider = provider
		}
	}
	
	// If no default provider is specified, use the first one
	if s.defaultProvider == nil {
		s.defaultProvider = config.Providers[0]
	}
	
	// Start background worker for async sending
	go s.processQueue()
	
	return s, nil
}

// Send sends an email immediately
func (s *emailService) Send(ctx context.Context, message *Message) error {
	if err := s.validateMessage(message); err != nil {
		return err
	}
	
	s.logger.Debugf("Sending email to %v with subject: %s", message.To, message.Subject)
	
	// Set default From if not provided
	if message.From == "" {
		message.From = s.config.DefaultFrom
	}
	
	return s.defaultProvider.Send(ctx, message)
}

// SendAsync queues an email for async delivery with retry
func (s *emailService) SendAsync(ctx context.Context, message *Message) (string, error) {
	if err := s.validateMessage(message); err != nil {
		return "", err
	}
	
	// Set default From if not provided
	if message.From == "" {
		message.From = s.config.DefaultFrom
	}
	
	id := uuid.New().String()
	
	qm := &queuedMessage{
		id:       id,
		message:  message,
		provider: s.defaultProvider,
		status:   StatusPending,
	}
	
	s.mu.Lock()
	s.queue[id] = qm
	s.mu.Unlock()
	
	s.logger.Debugf("Queued email %s to %v with subject: %s", id, message.To, message.Subject)
	
	return id, nil
}

// GetStatus retrieves the status of an async email
func (s *emailService) GetStatus(ctx context.Context, id string) (MessageStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	qm, exists := s.queue[id]
	if !exists {
		return "", fmt.Errorf("email with ID %s not found", id)
	}
	
	return qm.status, nil
}

// SendTemplate sends an email using a template
func (s *emailService) SendTemplate(ctx context.Context, templateName string, data map[string]interface{}, message *Message) error {
	// Template rendering would go here
	// For now, we'll just send the plain message
	return s.Send(ctx, message)
}

// validateMessage validates an email message
func (s *emailService) validateMessage(message *Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}
	
	if len(message.To) == 0 {
		return errors.New("at least one recipient is required")
	}
	
	if message.Subject == "" {
		return errors.New("subject is required")
	}
	
	if message.PlainText == "" && message.HtmlBody == "" {
		return errors.New("either plain text or HTML body is required")
	}
	
	return nil
}

// processQueue processes the queue of emails
func (s *emailService) processQueue() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	
	for range ticker.C {
		s.processQueueOnce()
	}
}

// processQueueOnce processes the queue once
func (s *emailService) processQueueOnce() {
	s.mu.Lock()
	
	// Get all pending messages
	var pending []*queuedMessage
	for _, qm := range s.queue {
		if qm.status == StatusPending && 
		   (qm.retries == 0 || time.Since(qm.lastAttempt) >= s.config.RetryInterval) {
			pending = append(pending, qm)
		}
	}
	
	s.mu.Unlock()
	
	// Process pending messages
	for _, qm := range pending {
		s.mu.Lock()
		
		// Check if message still exists and is still pending
		currentQm, exists := s.queue[qm.id]
		if !exists || currentQm.status != StatusPending {
			s.mu.Unlock()
			continue
		}
		
		qm.lastAttempt = time.Now()
		qm.retries++
		
		s.mu.Unlock()
		
		// Send the message
		err := qm.provider.Send(context.Background(), qm.message)
		
		s.mu.Lock()
		
		// Update status
		if err != nil {
			s.logger.Errorf("Failed to send email %s: %v", qm.id, err)
			qm.errorMsg = err.Error()
			
			if qm.retries >= s.config.MaxRetries {
				qm.status = StatusFailed
				s.logger.Errorf("Email %s failed after %d retries", qm.id, qm.retries)
			}
		} else {
			qm.status = StatusSent
			s.logger.Infof("Email %s sent successfully", qm.id)
		}
		
		s.mu.Unlock()
	}
	
	// Clean up old messages
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Remove messages older than a day
	cutoff := time.Now().Add(-24 * time.Hour)
	for id, qm := range s.queue {
		if (qm.status == StatusSent || qm.status == StatusFailed) && qm.lastAttempt.Before(cutoff) {
			delete(s.queue, id)
		}
	}
}