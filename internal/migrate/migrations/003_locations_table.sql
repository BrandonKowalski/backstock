CREATE TABLE locations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO locations (name) VALUES ('fridge'), ('inside_freezer'), ('garage_freezer'), ('pantry');

-- Drop the CHECK constraint by recreating the stock table
CREATE TABLE stock_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    location TEXT NOT NULL,
    quantity REAL NOT NULL DEFAULT 0,
    date_added TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

INSERT INTO stock_new SELECT * FROM stock;
DROP TABLE stock;
ALTER TABLE stock_new RENAME TO stock;

UPDATE schema_version SET version = 3;
