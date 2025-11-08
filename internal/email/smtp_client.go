package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/brandon/mcp-email/internal/config"
	"github.com/sirupsen/logrus"
)

// SMTPClient wraps an SMTP client
type SMTPClient struct {
	config *config.AccountConfig
	logger *logrus.Logger
}

// EmailMessage represents an email to be sent
type EmailMessage struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	BodyText    string
	BodyHTML    string
	Attachments []Attachment
	ReplyTo     string
	InReplyTo   string
}

// Attachment represents an email attachment
type Attachment struct {
	Filename string
	Content  []byte
	MimeType string
}

// NewSMTPClient creates a new SMTP client
func NewSMTPClient(cfg *config.AccountConfig) (*SMTPClient, error) {
	return &SMTPClient{
		config: cfg,
		logger: logrus.New(),
	}, nil
}

// Send sends an email
func (c *SMTPClient) Send(msg *EmailMessage) error {
	// Create message
	emailBytes, err := c.createMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Connect to server
	addr := fmt.Sprintf("%s:%d", c.config.SMTPHost, c.config.SMTPPort)
	
	// Determine if TLS is needed
	useTLS := c.config.SMTPPort == 465
	
	var auth smtp.Auth
	if c.config.SMTPPassword != "" {
		auth = smtp.PlainAuth("", c.config.SMTPUsername, c.config.SMTPPassword, c.config.SMTPHost)
	}

	if useTLS {
		// TLS connection (port 465)
		conn, err := tls.Dial("tcp", addr, &tls.Config{
			ServerName: c.config.SMTPHost,
		})
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, c.config.SMTPHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Close()

		// Auth
		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("failed to authenticate: %w", err)
			}
		}

		// Set sender
		if err := client.Mail(c.config.SMTPUsername); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		// Set recipients
		recipients := append(append(msg.To, msg.Cc...), msg.Bcc...)
		for _, to := range recipients {
			if err := client.Rcpt(to); err != nil {
				return fmt.Errorf("failed to set recipient %s: %w", to, err)
			}
		}

		// Send data
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to send data command: %w", err)
		}

		if _, err := w.Write(emailBytes); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}

		return client.Quit()
	} else {
		// StartTLS connection (port 587)
		client, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer client.Close()

		// Start TLS
		if err := client.StartTLS(&tls.Config{
			ServerName: c.config.SMTPHost,
		}); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}

		// Auth
		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("failed to authenticate: %w", err)
			}
		}

		// Set sender
		if err := client.Mail(c.config.SMTPUsername); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		// Set recipients
		recipients := append(append(msg.To, msg.Cc...), msg.Bcc...)
		for _, to := range recipients {
			if err := client.Rcpt(to); err != nil {
				return fmt.Errorf("failed to set recipient %s: %w", to, err)
			}
		}

		// Send data
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to send data command: %w", err)
		}

		if _, err := w.Write(emailBytes); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}

		return client.Quit()
	}
}

// createMessage creates an email message in MIME format
func (c *SMTPClient) createMessage(msg *EmailMessage) ([]byte, error) {
	var buf bytes.Buffer

	// Write headers manually (simpler approach)
	buf.WriteString(fmt.Sprintf("From: %s\r\n", c.config.SMTPUsername))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	if len(msg.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(msg.Cc, ", ")))
	}
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	if msg.ReplyTo != "" {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", msg.ReplyTo))
	}
	if msg.InReplyTo != "" {
		buf.WriteString(fmt.Sprintf("In-Reply-To: %s\r\n", msg.InReplyTo))
	}

	// Set content type
	if msg.BodyHTML != "" {
		buf.WriteString("Content-Type: text/html; charset=utf-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(msg.BodyHTML)
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(msg.BodyText)
	}

	return buf.Bytes(), nil
}

// SetLogger sets the logger for the client
func (c *SMTPClient) SetLogger(logger *logrus.Logger) {
	c.logger = logger
}

