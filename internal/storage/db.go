package storage

import (
	"database/sql"
	_ "embed"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed migrations/001_init.sql
var initSQL string

type DB struct {
	conn *sql.DB
}

func NewDB(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", path+"?_busy_timeout=5000&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)

	if _, err := conn.Exec(initSQL); err != nil {
		conn.Close()
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Conn() *sql.DB {
	return db.conn
}
