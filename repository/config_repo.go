package repository

import (
	"database/sql"
	"fmt"
)

// ConfigRepository provides key-value access to the config table.
type ConfigRepository struct {
	db *sql.DB
}

// NewConfigRepository creates a new ConfigRepository.
func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// Get retrieves a config value by key. Returns empty string if not found.
func (r *ConfigRepository) Get(key string) (string, error) {
	var value string
	err := r.db.QueryRow(`SELECT value FROM config WHERE key = ?`, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("get config %q: %w", key, err)
	}
	return value, nil
}

// Set inserts or updates a config value.
func (r *ConfigRepository) Set(key, value string) error {
	_, err := r.db.Exec(
		`INSERT INTO config (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	if err != nil {
		return fmt.Errorf("set config %q: %w", key, err)
	}
	return nil
}
