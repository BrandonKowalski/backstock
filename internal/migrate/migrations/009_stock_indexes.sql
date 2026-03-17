CREATE INDEX IF NOT EXISTS idx_stock_item_id ON stock(item_id);
CREATE INDEX IF NOT EXISTS idx_stock_location ON stock(location);

UPDATE schema_version SET version = 9;
