package store

import "backstock/internal/model"

func (s *Store) ListLocations() ([]model.Location, error) {
	rows, err := s.db.Query("SELECT id, name, parent_id, is_food, exclude_default FROM locations ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []model.Location
	for rows.Next() {
		var loc model.Location
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.ParentID, &loc.IsFood, &loc.ExcludeDefault); err != nil {
			return nil, err
		}
		all = append(all, loc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Build tree: top-level locations with children nested
	byID := make(map[int]*model.Location)
	for i := range all {
		byID[all[i].ID] = &all[i]
	}

	var roots []model.Location
	for i := range all {
		if all[i].ParentID != nil {
			parent := byID[*all[i].ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, all[i])
			}
		} else {
			roots = append(roots, all[i])
		}
	}

	// Re-attach children to roots (since roots are copies)
	for i := range roots {
		if orig, ok := byID[roots[i].ID]; ok {
			roots[i].Children = orig.Children
		}
	}

	return roots, nil
}

// ListLocationNames returns a flat list of all location names (for stock queries).
func (s *Store) ListLocationNames() ([]string, error) {
	rows, err := s.db.Query("SELECT name FROM locations ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

func (s *Store) CreateLocation(loc *model.Location) error {
	res, err := s.db.Exec(
		"INSERT INTO locations (name, parent_id, is_food, exclude_default) VALUES (?, ?, ?, ?)",
		loc.Name, loc.ParentID, loc.IsFood, loc.ExcludeDefault,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	loc.ID = int(id)
	return nil
}

func (s *Store) UpdateLocation(loc *model.Location) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get old name to update stock references
	var oldName string
	if err := tx.QueryRow("SELECT name FROM locations WHERE id = ?", loc.ID).Scan(&oldName); err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE locations SET name = ?, parent_id = ?, is_food = ?, exclude_default = ? WHERE id = ?",
		loc.Name, loc.ParentID, loc.IsFood, loc.ExcludeDefault, loc.ID,
	)
	if err != nil {
		return err
	}

	// Update stock rows that reference the old name
	if oldName != loc.Name {
		if _, err = tx.Exec("UPDATE stock SET location = ? WHERE location = ?", loc.Name, oldName); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) DeleteLocation(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Move children to no parent before deleting
	if _, err := tx.Exec("UPDATE locations SET parent_id = NULL WHERE parent_id = ?", id); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM locations WHERE id = ?", id); err != nil {
		return err
	}

	return tx.Commit()
}
