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
	conn.SetMaxIdleConns(1)

	if _, err := conn.Exec(initSQL); err != nil {
		conn.Close()
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Conn() *sql.DB {
	return db.conn
}

func (db *DB) migrate() error {
	return db.migrateHistory()
}

func (db *DB) migrateHistory() error {
	histCols, err := db.getColumns("clean_history")
	if err != nil {
		return err
	}
	if !histCols["failed_count"] {
		if _, err := db.conn.Exec("ALTER TABLE clean_history ADD COLUMN failed_count INTEGER DEFAULT 0"); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) getColumns(table string) (map[string]bool, error) {
	columns := make(map[string]bool)
	rows, err := db.conn.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name string
		var typ string
		var notnull int
		var dfltValue interface{}
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			continue
		}
		columns[name] = true
	}
	return columns, nil
}
