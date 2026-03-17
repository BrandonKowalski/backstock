ALTER TABLE categories ADD COLUMN is_food INTEGER NOT NULL DEFAULT 1;

UPDATE schema_version SET version = 5;
