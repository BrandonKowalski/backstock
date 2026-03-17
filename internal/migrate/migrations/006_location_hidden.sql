ALTER TABLE locations ADD COLUMN exclude_default INTEGER NOT NULL DEFAULT 0;

UPDATE schema_version SET version = 6;
