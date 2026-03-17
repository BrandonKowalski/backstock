package store

import (
	"backstock/internal/model"
)

func (s *Store) ListUnits() ([]model.Unit, error) {
	rows, err := s.db.Query("SELECT id, name, abbreviation, is_food FROM units ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []model.Unit
	for rows.Next() {
		var u model.Unit
		if err := rows.Scan(&u.ID, &u.Name, &u.Abbreviation, &u.IsFood); err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

func (s *Store) CreateUnit(u *model.Unit) error {
	res, err := s.db.Exec("INSERT INTO units (name, abbreviation, is_food) VALUES (?, ?, ?)", u.Name, u.Abbreviation, u.IsFood)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

func (s *Store) UpdateUnit(u *model.Unit) error {
	_, err := s.db.Exec("UPDATE units SET name = ?, abbreviation = ?, is_food = ? WHERE id = ?", u.Name, u.Abbreviation, u.IsFood, u.ID)
	return err
}

func (s *Store) DeleteUnit(id int) error {
	_, err := s.db.Exec("DELETE FROM units WHERE id = ?", id)
	return err
}
