PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS clean_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    scan_level INTEGER NOT NULL,
    total_files INTEGER NOT NULL,
    total_size INTEGER NOT NULL,
    freed_size INTEGER NOT NULL,
    failed_count INTEGER DEFAULT 0,
    status TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clean_details (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    history_id INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    rule_id TEXT NOT NULL,
    risk_score INTEGER NOT NULL,
    action TEXT NOT NULL,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (history_id) REFERENCES clean_history(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_clean_details_history ON clean_details(history_id);
CREATE INDEX IF NOT EXISTS idx_clean_details_rule ON clean_details(rule_id);

CREATE TABLE IF NOT EXISTS user_config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rule_status (
    rule_id TEXT PRIMARY KEY,
    enabled BOOLEAN NOT NULL DEFAULT 1,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
