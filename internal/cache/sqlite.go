package cache

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

// Cache represents the SQLite cache
type Cache struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewCache creates a new cache instance
func NewCache(dbPath string, logger *logrus.Logger) (*Cache, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	cache := &Cache{
		db:     db,
		logger: logger,
	}

	// Initialize schema
	if err := cache.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	logger.WithField("path", dbPath).Info("Cache initialized")
	return cache, nil
}

// initSchema initializes the database schema
func (c *Cache) initSchema() error {
	if _, err := c.db.Exec(Schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}

// Close closes the database connection
func (c *Cache) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// DB returns the underlying database connection (for use in store.go)
func (c *Cache) DB() *sql.DB {
	return c.db
}
