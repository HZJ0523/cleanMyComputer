package storage

import (
	"database/sql"
	"errors"

	"github.com/hzj0523/cleanMyComputer/internal/models"
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

func (c *Config) GetConfig() *models.Config {
	cfg := &models.Config{
		OldFileDays: 30,
		ScanWorkers: 4,
		Language:    "zh-CN",
	}
	if v, err := c.Get("auto_clean_enabled"); err == nil && v == "true" {
		cfg.AutoCleanEnabled = true
	}
	if v, err := c.Get("old_file_days"); err == nil && v != "" {
		if n, e := parseInt(v); e == nil && n > 0 {
			cfg.OldFileDays = n
		}
	}
	if v, err := c.Get("scan_workers"); err == nil && v != "" {
		if n, e := parseInt(v); e == nil && n > 0 {
			cfg.ScanWorkers = n
		}
	}
	if v, err := c.Get("language"); err == nil && v != "" {
		cfg.Language = v
	}
	return cfg
}

func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("invalid")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
