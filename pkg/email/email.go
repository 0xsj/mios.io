package email

import "context"

type Service interface {
	// Send sends an email immediately
	Send(ctx context.Context, message *Message) error
	
	// SendAsync queues an email for async delivery with retry
	SendAsync(ctx context.Context, message *Message) (string, error)
	
	// GetStatus retrieves the status of an async email
	GetStatus(ctx context.Context, id string) (MessageStatus, error)
	
	// SendTemplate sends an email using a template
	SendTemplate(ctx context.Context, templateName string, data map[string]interface{}, message *Message) error
}


type Message struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	PlainText   string
	HtmlBody    string
	Attachments []Attachment
	Headers     map[string]string
}

// Attachment represents a file attachment
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

type MessageStatus string

const (
	StatusPending    MessageStatus = "pending"
	StatusSent       MessageStatus = "sent"
	StatusFailed     MessageStatus = "failed"
	StatusDelivered  MessageStatus = "delivered"
)

type Provider interface {
	Send(ctx context.Context, message *Message) error
	Name() string
}