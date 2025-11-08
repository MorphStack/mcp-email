package tools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/brandon/mcp-email/internal/cache"
	"github.com/brandon/mcp-email/internal/config"
	"github.com/brandon/mcp-email/internal/email"
)

// GetEmailTool retrieves a full email by ID
type GetEmailTool struct {
	config       *config.Config
	emailManager *email.Manager
	cacheStore   *cache.Store
	logger       *logrus.Logger
}

// NewGetEmailTool creates a new get email tool
func NewGetEmailTool(cfg *config.Config, emailManager *email.Manager, cacheStore *cache.Store, logger *logrus.Logger) *GetEmailTool {
	return &GetEmailTool{
		config:       cfg,
		emailManager: emailManager,
		cacheStore:   cacheStore,
		logger:       logger,
	}
}

// Name returns the tool name
func (t *GetEmailTool) Name() string {
	return "get_email"
}

// Description returns the tool description
func (t *GetEmailTool) Description() string {
	return "Retrieve full email by ID from cache or IMAP"
}

// InputSchema returns the JSON schema for tool inputs
func (t *GetEmailTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"email_id": map[string]interface{}{
				"type":        "integer",
				"description": "Email ID (from search results)",
			},
			"account_name": map[string]interface{}{
				"type":        "string",
				"description": "Optional: Account name if needed",
			},
		},
		"required": []string{"email_id"},
	}
}

// Execute executes the tool
func (t *GetEmailTool) Execute(params map[string]interface{}) (interface{}, error) {
	// Parse email_id
	var emailID int64
	if id, ok := params["email_id"].(float64); ok {
		emailID = int64(id)
	} else if idStr, ok := params["email_id"].(string); ok {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid email_id: %w", err)
		}
		emailID = id
	} else {
		return nil, fmt.Errorf("email_id is required")
	}

	// Get email from cache
	email, err := t.cacheStore.GetEmail(emailID)
	if err != nil {
		return nil, fmt.Errorf("failed to get email: %w", err)
	}

	// Convert to JSON-serializable format
	result := map[string]interface{}{
		"id":           email.ID,
		"account_id":   email.AccountID,
		"account_name": email.AccountName,
		"folder_id":    email.FolderID,
		"folder_path":  email.FolderPath,
		"uid":          email.UID,
		"message_id":   email.MessageID,
		"subject":      email.Subject,
		"sender_name":  email.SenderName,
		"sender_email": email.SenderEmail,
		"recipients":   email.Recipients,
		"date":         email.Date.Format(time.RFC3339),
		"body_text":    email.BodyText,
		"body_html":    email.BodyHTML,
		"headers":      email.Headers,
		"flags":        email.Flags,
		"cached_at":    email.CachedAt.Format(time.RFC3339),
	}

	return result, nil
}
