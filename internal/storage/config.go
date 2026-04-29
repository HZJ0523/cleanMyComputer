package storage

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("config key not found")

type Config struct {
	db *DB
}

func NewConfig(db *DB) *Config {
	return &Config{db: db}
}

func (c *Config) Get(key string) (string, error) {
	var value string
	err := c.db.Conn().QueryRow("SELECT value FROM user_config WHERE key = ?", key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return value, err
}

func (c *Config) Set(key, value string) error {
	_, err := c.db.Conn().Exec(`
		INSERT INTO user_config (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
	`, key, value, value)
	return err
}
