// pkg/email/client.go
package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"

	"github.com/0xsj/mios.io/log"
)

// EmailClient handles sending emails
type EmailClient struct {
	host            string
	port            int
	from            string
	fromName        string
	logger          log.Logger
	templateManager *TemplateManager // Add template manager
}

// Message represents an email message
type Message struct {
	To      []string
	Subject string
	Body    string // Used when template is not specified
	IsHTML  bool
}

// EmailData represents data for email template rendering
type EmailData struct {
	Username   string
	Link       string
	AppName    string
	Year       int
	CustomData map[string]string
}

// NewEmailClient creates a new email client from environment variables
func NewEmailClient(logger log.Logger, templateManager *TemplateManager) *EmailClient {
	port, _ := strconv.Atoi(getEnvOrDefault("EMAIL_PORT", "1025"))

	return &EmailClient{
		host:            getEnvOrDefault("EMAIL_HOST", "localhost"),
		port:            port,
		from:            getEnvOrDefault("EMAIL_FROM", "noreply@example.com"),
		fromName:        getEnvOrDefault("EMAIL_FROM_NAME", "Your App"),
		logger:          logger.WithLayer("EmailClient"),
		templateManager: templateManager,
	}
}

// Send sends a simple email with provided body
func (c *EmailClient) Send(msg Message) error {
	c.logger.Debugf("Sending email to %v: %s", msg.To, msg.Subject)

	return c.sendEmail(msg.To, msg.Subject, msg.Body, msg.IsHTML)
}

// SendTemplate sends an email using a template
func (c *EmailClient) SendTemplate(to []string, subject, templateName string, data interface{}) error {
	c.logger.Debugf("Sending template email '%s' to %v", templateName, to)

	// Render the template
	body, err := c.templateManager.Render(templateName, data)
	if err != nil {
		c.logger.Errorf("Failed to render template %s: %v", templateName, err)
		return err
	}

	// Determine if it's HTML based on template extension
	isHTML := true
	if len(templateName) > 4 && templateName[len(templateName)-4:] == ".txt" {
		isHTML = false
	}

	return c.sendEmail(to, subject, body, isHTML)
}

// sendEmail is an internal method that handles the actual sending
func (c *EmailClient) sendEmail(to []string, subject, body string, isHTML bool) error {
	// Build email
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	// Set sender with name if available
	from := c.from
	if c.fromName != "" {
		from = fmt.Sprintf("%s <%s>", c.fromName, c.from)
	}

	// Build email headers and body
	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}

	emailContent := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: %s; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		from,
		to[0], // For simplicity, just use the first recipient in headers
		subject,
		contentType,
		body,
	)

	// Send email without authentication (for MailHog)
	err := smtp.SendMail(addr, nil, c.from, to, []byte(emailContent))
	if err != nil {
		c.logger.Errorf("Failed to send email: %v", err)
		return err
	}

	c.logger.Infof("Email sent successfully to %v", to)
	return nil
}

// Helper to get environment variables with defaults
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
