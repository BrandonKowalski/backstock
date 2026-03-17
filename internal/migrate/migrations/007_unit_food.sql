ALTER TABLE units ADD COLUMN is_food INTEGER;

UPDATE schema_version SET version = 7;
