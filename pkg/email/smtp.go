// pkg/email/smtp.go
package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

// SMTPProvider implements the Provider interface for SMTP
type SMTPProvider struct {
	config SMTPConfig
}

// NewSMTPProvider creates a new SMTP provider
func NewSMTPProvider(config SMTPConfig) *SMTPProvider {
	return &SMTPProvider{
		config: config,
	}
}

// Name returns the provider name
func (p *SMTPProvider) Name() string {
	return "smtp"
}

// Send sends an email via SMTP
func (p *SMTPProvider) Send(ctx context.Context, message *Message) error {
	addr := fmt.Sprintf("%s:%d", p.config.Host, p.config.Port)
	auth := smtp.PlainAuth("", p.config.Username, p.config.Password, p.config.Host)
	
	from := message.From
	if from == "" {
		from = p.config.From
	}
	
	// Build email headers
	headers := make(map[string]string)
	if message.Headers != nil {
		for k, v := range message.Headers {
			headers[k] = v
		}
	}
	
	headers["From"] = from
	headers["To"] = strings.Join(message.To, ", ")
	if len(message.Cc) > 0 {
		headers["Cc"] = strings.Join(message.Cc, ", ")
	}
	headers["Subject"] = message.Subject
	headers["MIME-Version"] = "1.0"
	
	// Build email body
	var body strings.Builder
	
	// Add headers
	for k, v := range headers {
		body.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	
	// Add email content
	if message.HtmlBody != "" {
		body.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		body.WriteString(message.HtmlBody)
	} else {
		body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		body.WriteString(message.PlainText)
	}
	
	// Get all recipients
	recipients := append([]string{}, message.To...)
	recipients = append(recipients, message.Cc...)
	recipients = append(recipients, message.Bcc...)
	
	// Send the email
	if p.config.UseTLS {
		tlsConfig := &tls.Config{
			ServerName: p.config.Host,
		}
		
		client, err := smtp.Dial(addr)
		if err != nil {
			return err
		}
		defer client.Close()
		
		if err = client.StartTLS(tlsConfig); err != nil {
			return err
		}
		
		if err = client.Auth(auth); err != nil {
			return err
		}
		
		if err = client.Mail(from); err != nil {
			return err
		}
		
		for _, recipient := range recipients {
			if err = client.Rcpt(recipient); err != nil {
				return err
			}
		}
		
		w, err := client.Data()
		if err != nil {
			return err
		}
		
		_, err = w.Write([]byte(body.String()))
		if err != nil {
			return err
		}
		
		err = w.Close()
		if err != nil {
			return err
		}
		
		return client.Quit()
	} else {
		return smtp.SendMail(addr, auth, from, recipients, []byte(body.String()))
	}
}