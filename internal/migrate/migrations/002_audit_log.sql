CREATE TABLE audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_name TEXT NOT NULL,
    quantity REAL NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

UPDATE schema_version SET version = 2;
