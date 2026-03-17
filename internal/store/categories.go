package store

import (
	"backstock/internal/model"
)

func (s *Store) ListCategories() ([]model.Category, error) {
	rows, err := s.db.Query("SELECT id, name, is_food FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.IsFood); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (s *Store) CreateCategory(c *model.Category) error {
	res, err := s.db.Exec("INSERT INTO categories (name, is_food) VALUES (?, ?)", c.Name, c.IsFood)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	return nil
}

func (s *Store) UpdateCategory(c *model.Category) error {
	_, err := s.db.Exec("UPDATE categories SET name = ?, is_food = ? WHERE id = ?", c.Name, c.IsFood, c.ID)
	return err
}

func (s *Store) DeleteCategory(id int) error {
	_, err := s.db.Exec("DELETE FROM categories WHERE id = ?", id)
	return err
}
