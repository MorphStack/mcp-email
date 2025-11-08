package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// Cache settings
	CachePath          string
	SearchResultLimit  int
	LogLevel           string

	// Accounts
	Accounts []AccountConfig
}

// AccountConfig holds configuration for a single email account
type AccountConfig struct {
	Name string

	// IMAP settings
	IMAPHost     string
	IMAPPort     int
	IMAPUsername string
	IMAPPassword string

	// SMTP settings
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		CachePath:         getEnv("CACHE_PATH", "/data/email_cache.db"),
		SearchResultLimit: getEnvInt("SEARCH_RESULT_LIMIT", 100),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
	}

	// Load accounts
	accounts, err := loadAccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to load accounts: %w", err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("no email accounts configured")
	}

	cfg.Accounts = accounts
	return cfg, nil
}

// loadAccounts loads email account configurations from environment variables
func loadAccounts() ([]AccountConfig, error) {
	var accounts []AccountConfig

	// First, try single account configuration (for backward compatibility)
	if hasSingleAccount() {
		account, err := loadSingleAccount()
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, *account)
		return accounts, nil
	}

	// Load multiple accounts (ACCOUNT_1_*, ACCOUNT_2_*, etc.)
	accountNum := 1
	for {
		account, err := loadAccountByNumber(accountNum)
		if err != nil {
			break // No more accounts
		}
		accounts = append(accounts, *account)
		accountNum++
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found in environment variables")
	}

	return accounts, nil
}

// hasSingleAccount checks if single account configuration exists
func hasSingleAccount() bool {
	return getEnv("IMAP_HOST", "") != "" && getEnv("SMTP_HOST", "") != ""
}

// loadSingleAccount loads a single account from environment variables
func loadSingleAccount() (*AccountConfig, error) {
	imapHost := getEnv("IMAP_HOST", "")
	imapPort := getEnvInt("IMAP_PORT", 993)
	imapUsername := getEnv("IMAP_USERNAME", "")
	imapPassword := getEnv("IMAP_PASSWORD", "")

	smtpHost := getEnv("SMTP_HOST", "")
	smtpPort := getEnvInt("SMTP_PORT", 587)
	smtpUsername := getEnv("SMTP_USERNAME", "")
	smtpPassword := getEnv("SMTP_PASSWORD", "")

	if imapHost == "" || smtpHost == "" {
		return nil, fmt.Errorf("IMAP_HOST and SMTP_HOST are required")
	}

	if imapUsername == "" || smtpUsername == "" {
		return nil, fmt.Errorf("IMAP_USERNAME and SMTP_USERNAME are required")
	}

	if imapPassword == "" || smtpPassword == "" {
		return nil, fmt.Errorf("IMAP_PASSWORD and SMTP_PASSWORD are required")
	}

	// Default account name
	name := getEnv("ACCOUNT_NAME", "default")
	if name == "" {
		name = "default"
	}

	return &AccountConfig{
		Name:         name,
		IMAPHost:     imapHost,
		IMAPPort:     imapPort,
		IMAPUsername: imapUsername,
		IMAPPassword: imapPassword,
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
	}, nil
}

// loadAccountByNumber loads an account by number (ACCOUNT_1_*, ACCOUNT_2_*, etc.)
func loadAccountByNumber(num int) (*AccountConfig, error) {
	prefix := fmt.Sprintf("ACCOUNT_%d_", num)

	name := getEnv(prefix+"NAME", "")
	if name == "" {
		return nil, fmt.Errorf("account %d: NAME is required", num)
	}

	imapHost := getEnv(prefix+"IMAP_HOST", "")
	imapPort := getEnvInt(prefix+"IMAP_PORT", 993)
	imapUsername := getEnv(prefix+"IMAP_USERNAME", "")
	imapPassword := getEnv(prefix+"IMAP_PASSWORD", "")

	smtpHost := getEnv(prefix+"SMTP_HOST", "")
	smtpPort := getEnvInt(prefix+"SMTP_PORT", 587)
	smtpUsername := getEnv(prefix+"SMTP_USERNAME", "")
	smtpPassword := getEnv(prefix+"SMTP_PASSWORD", "")

	if imapHost == "" || smtpHost == "" {
		return nil, fmt.Errorf("account %d: IMAP_HOST and SMTP_HOST are required", num)
	}

	if imapUsername == "" || smtpUsername == "" {
		return nil, fmt.Errorf("account %d: IMAP_USERNAME and SMTP_USERNAME are required", num)
	}

	if imapPassword == "" || smtpPassword == "" {
		return nil, fmt.Errorf("account %d: IMAP_PASSWORD and SMTP_PASSWORD are required", num)
	}

	return &AccountConfig{
		Name:         name,
		IMAPHost:     imapHost,
		IMAPPort:     imapPort,
		IMAPUsername: imapUsername,
		IMAPPassword: imapPassword,
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as an integer or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetAccountByName finds an account by name
func (c *Config) GetAccountByName(name string) (*AccountConfig, error) {
	for i := range c.Accounts {
		if c.Accounts[i].Name == name {
			return &c.Accounts[i], nil
		}
	}
	return nil, fmt.Errorf("account not found: %s", name)
}

// GetDefaultAccount returns the first account (or default account if named "default")
func (c *Config) GetDefaultAccount() *AccountConfig {
	if len(c.Accounts) == 0 {
		return nil
	}

	// Try to find "default" account first
	for i := range c.Accounts {
		if c.Accounts[i].Name == "default" {
			return &c.Accounts[i]
		}
	}

	// Return first account
	return &c.Accounts[0]
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.CachePath == "" {
		return fmt.Errorf("CACHE_PATH is required")
	}

	if c.SearchResultLimit < 1 || c.SearchResultLimit > 1000 {
		return fmt.Errorf("SEARCH_RESULT_LIMIT must be between 1 and 1000")
	}

	if len(c.Accounts) == 0 {
		return fmt.Errorf("at least one account must be configured")
	}

	// Validate each account
	for i := range c.Accounts {
		acc := &c.Accounts[i]
		if acc.IMAPHost == "" {
			return fmt.Errorf("account %s: IMAP_HOST is required", acc.Name)
		}
		if acc.SMTPHost == "" {
			return fmt.Errorf("account %s: SMTP_HOST is required", acc.Name)
		}
		if acc.IMAPPort < 1 || acc.IMAPPort > 65535 {
			return fmt.Errorf("account %s: invalid IMAP_PORT", acc.Name)
		}
		if acc.SMTPPort < 1 || acc.SMTPPort > 65535 {
			return fmt.Errorf("account %s: invalid SMTP_PORT", acc.Name)
		}
	}

	return nil
}

// AccountNames returns a list of all account names
func (c *Config) AccountNames() []string {
	names := make([]string, len(c.Accounts))
	for i := range c.Accounts {
		names[i] = c.Accounts[i].Name
	}
	return names
}

