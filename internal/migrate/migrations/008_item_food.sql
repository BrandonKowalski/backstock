ALTER TABLE items ADD COLUMN is_food INTEGER NOT NULL DEFAULT 1;

-- Set is_food based on stock locations for existing items
UPDATE items SET is_food = 0 WHERE id IN (
    SELECT DISTINCT s.item_id FROM stock s
    JOIN locations l ON l.name = s.location
    WHERE l.is_food = 0
) AND id NOT IN (
    SELECT DISTINCT s.item_id FROM stock s
    JOIN locations l ON l.name = s.location
    WHERE l.is_food = 1
);

UPDATE schema_version SET version = 8;
