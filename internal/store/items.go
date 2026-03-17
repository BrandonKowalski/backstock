package store

import (
	"fmt"
	"strings"
	"backstock/internal/model"
)

func (s *Store) ListItems(f model.ItemFilter) ([]model.Item, error) {
	query := `
		SELECT DISTINCT i.id, i.name, i.is_food, i.unit_id, i.package_size, i.expiration_date,
			i.best_by_date, i.low_quantity_threshold, i.created_at, i.updated_at,
			COALESCE(u.id, 0), COALESCE(u.name, ''), COALESCE(u.abbreviation, ''),
			COALESCE(sq.total, 0)
		FROM items i
		LEFT JOIN units u ON i.unit_id = u.id
		LEFT JOIN (SELECT item_id, SUM(quantity) as total FROM stock GROUP BY item_id) sq ON sq.item_id = i.id`

	var where []string
	var args []any

	if f.Location != "" {
		query += " JOIN stock s ON s.item_id = i.id"
		where = append(where, "s.location IN (SELECT name FROM locations WHERE name = ? OR parent_id = (SELECT id FROM locations WHERE name = ?))")
		args = append(args, f.Location, f.Location)
	}
	if f.Category != "" {
		query += " JOIN item_categories ic ON ic.item_id = i.id JOIN categories c ON c.id = ic.category_id"
		where = append(where, "c.name = ?")
		args = append(args, f.Category)
	}
	if f.Search != "" {
		where = append(where, "i.name LIKE ?")
		args = append(args, "%"+f.Search+"%")
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	switch f.Sort {
	case "name":
		query += " ORDER BY i.name"
	case "expiration":
		query += " ORDER BY COALESCE(i.expiration_date, '9999-12-31')"
	case "recent":
		query += " ORDER BY i.created_at DESC"
	default:
		query += " ORDER BY i.name"
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer rows.Close()

	var items []model.Item
	var itemIDs []any
	for rows.Next() {
		var it model.Item
		var u model.Unit
		if err := rows.Scan(&it.ID, &it.Name, &it.IsFood, &it.UnitID, &it.PackageSize, &it.ExpirationDate,
			&it.BestByDate, &it.LowQuantityThreshold, &it.CreatedAt, &it.UpdatedAt,
			&u.ID, &u.Name, &u.Abbreviation, &it.TotalQuantity); err != nil {
			return nil, err
		}
		if u.ID != 0 {
			it.Unit = &u
		}
		items = append(items, it)
		itemIDs = append(itemIDs, it.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return items, nil
	}

	// Batch load categories for all items
	catMap, err := s.batchItemCategories(itemIDs)
	if err != nil {
		return nil, err
	}

	// Batch load stock for all items
	stockMap, err := s.batchListStock(itemIDs)
	if err != nil {
		return nil, err
	}

	for idx := range items {
		items[idx].Categories = catMap[items[idx].ID]
		items[idx].Stock = stockMap[items[idx].ID]
	}

	return items, nil
}

func (s *Store) GetItem(id int) (*model.Item, error) {
	var it model.Item
	var u model.Unit
	err := s.db.QueryRow(`
		SELECT i.id, i.name, i.is_food, i.unit_id, i.package_size, i.expiration_date,
			i.best_by_date, i.low_quantity_threshold, i.created_at, i.updated_at,
			COALESCE(u.id, 0), COALESCE(u.name, ''), COALESCE(u.abbreviation, ''),
			COALESCE(sq.total, 0)
		FROM items i
		LEFT JOIN units u ON i.unit_id = u.id
		LEFT JOIN (SELECT item_id, SUM(quantity) as total FROM stock GROUP BY item_id) sq ON sq.item_id = i.id
		WHERE i.id = ?`, id).Scan(&it.ID, &it.Name, &it.IsFood, &it.UnitID, &it.PackageSize, &it.ExpirationDate,
		&it.BestByDate, &it.LowQuantityThreshold, &it.CreatedAt, &it.UpdatedAt,
		&u.ID, &u.Name, &u.Abbreviation, &it.TotalQuantity)
	if err != nil {
		return nil, err
	}
	if u.ID != 0 {
		it.Unit = &u
	}

	cats, err := s.itemCategories(it.ID)
	if err != nil {
		return nil, err
	}
	it.Categories = cats

	stock, err := s.ListStock(it.ID)
	if err != nil {
		return nil, err
	}
	it.Stock = stock

	return &it, nil
}

func (s *Store) CreateItem(it *model.Item, categoryIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO items (name, is_food, unit_id, package_size, expiration_date, best_by_date, low_quantity_threshold)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		it.Name, it.IsFood, it.UnitID, it.PackageSize, it.ExpirationDate, it.BestByDate, it.LowQuantityThreshold)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	it.ID = int(id)

	for _, cid := range categoryIDs {
		if _, err := tx.Exec("INSERT INTO item_categories (item_id, category_id) VALUES (?, ?)", it.ID, cid); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) UpdateItem(it *model.Item, categoryIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE items SET name=?, is_food=?, unit_id=?, package_size=?, expiration_date=?, best_by_date=?,
			low_quantity_threshold=?, updated_at=datetime('now')
		WHERE id=?`,
		it.Name, it.IsFood, it.UnitID, it.PackageSize, it.ExpirationDate, it.BestByDate,
		it.LowQuantityThreshold, it.ID)
	if err != nil {
		return err
	}

	if categoryIDs != nil {
		if _, err = tx.Exec("DELETE FROM item_categories WHERE item_id = ?", it.ID); err != nil {
			return err
		}
		for _, cid := range categoryIDs {
			if _, err = tx.Exec("INSERT INTO item_categories (item_id, category_id) VALUES (?, ?)", it.ID, cid); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *Store) DeleteItem(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Explicit cleanup in case foreign_keys pragma is off
	if _, err := tx.Exec("DELETE FROM stock WHERE item_id = ?", id); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM item_categories WHERE item_id = ?", id); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM items WHERE id = ?", id); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) itemCategories(itemID int) ([]model.Category, error) {
	rows, err := s.db.Query(`
		SELECT c.id, c.name, c.is_food FROM categories c
		JOIN item_categories ic ON ic.category_id = c.id
		WHERE ic.item_id = ?
		ORDER BY c.name`, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.IsFood); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

// batchItemCategories loads categories for multiple items in a single query.
func (s *Store) batchItemCategories(itemIDs []any) (map[int][]model.Category, error) {
	placeholders := strings.Repeat("?,", len(itemIDs))
	placeholders = placeholders[:len(placeholders)-1]

	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT ic.item_id, c.id, c.name, c.is_food FROM categories c
		JOIN item_categories ic ON ic.category_id = c.id
		WHERE ic.item_id IN (%s)
		ORDER BY c.name`, placeholders), itemIDs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int][]model.Category)
	for rows.Next() {
		var itemID int
		var c model.Category
		if err := rows.Scan(&itemID, &c.ID, &c.Name, &c.IsFood); err != nil {
			return nil, err
		}
		result[itemID] = append(result[itemID], c)
	}
	return result, rows.Err()
}

// batchListStock loads stock for multiple items in a single query.
func (s *Store) batchListStock(itemIDs []any) (map[int][]model.Stock, error) {
	placeholders := strings.Repeat("?,", len(itemIDs))
	placeholders = placeholders[:len(placeholders)-1]

	rows, err := s.db.Query(fmt.Sprintf(
		"SELECT id, item_id, location, quantity, date_added, updated_at FROM stock WHERE item_id IN (%s) ORDER BY location",
		placeholders), itemIDs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int][]model.Stock)
	for rows.Next() {
		var st model.Stock
		if err := rows.Scan(&st.ID, &st.ItemID, &st.Location, &st.Quantity, &st.DateAdded, &st.UpdatedAt); err != nil {
			return nil, err
		}
		result[st.ItemID] = append(result[st.ItemID], st)
	}
	return result, rows.Err()
}
