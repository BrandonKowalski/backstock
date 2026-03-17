ALTER TABLE locations ADD COLUMN parent_id INTEGER REFERENCES locations(id);
ALTER TABLE locations ADD COLUMN is_food INTEGER NOT NULL DEFAULT 1;

UPDATE schema_version SET version = 4;
