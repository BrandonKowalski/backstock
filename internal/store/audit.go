package store

import "backstock/internal/model"

func (s *Store) AddAuditEntry(itemName string, quantity float64) error {
	_, err := s.db.Exec("INSERT INTO audit_log (item_name, quantity) VALUES (?, ?)", itemName, quantity)
	return err
}

func (s *Store) ListAuditLog() ([]model.AuditEntry, error) {
	rows, err := s.db.Query("SELECT id, item_name, quantity, created_at FROM audit_log ORDER BY created_at DESC LIMIT 200")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []model.AuditEntry
	for rows.Next() {
		var e model.AuditEntry
		if err := rows.Scan(&e.ID, &e.ItemName, &e.Quantity, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
