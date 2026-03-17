package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
)

//go:embed migrations/*.sql
var migrations embed.FS

const latestVersion = 9

type migration struct {
	version int
	file    string
}

var allMigrations = []migration{
	{1, "migrations/001_initial.sql"},
	{2, "migrations/002_audit_log.sql"},
	{3, "migrations/003_locations_table.sql"},
	{4, "migrations/004_location_parent.sql"},
	{5, "migrations/005_category_food.sql"},
	{6, "migrations/006_location_hidden.sql"},
	{7, "migrations/007_unit_food.sql"},
	{8, "migrations/008_item_food.sql"},
	{9, "migrations/009_stock_indexes.sql"},
}

func Run(db *sql.DB) error {
	var version int
	err := db.QueryRow("SELECT version FROM schema_version LIMIT 1").Scan(&version)
	if err != nil {
		version = 0
	}

	if version >= latestVersion {
		log.Printf("database at version %d, no migrations needed", version)
		return nil
	}

	for _, m := range allMigrations {
		if m.version <= version {
			continue
		}
		log.Printf("applying %s", m.file)
		data, err := migrations.ReadFile(m.file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", m.file, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", m.file, err)
		}
		if _, err := tx.Exec(string(data)); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", m.file, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", m.file, err)
		}
	}

	log.Printf("database migrated to version %d", latestVersion)
	return nil
}
