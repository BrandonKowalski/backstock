package store

import (
	"fmt"
	"backstock/internal/model"
)

func (s *Store) ListStock(itemID int) ([]model.Stock, error) {
	rows, err := s.db.Query(
		"SELECT id, item_id, location, quantity, date_added, updated_at FROM stock WHERE item_id = ? ORDER BY location",
		itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []model.Stock
	for rows.Next() {
		var st model.Stock
		if err := rows.Scan(&st.ID, &st.ItemID, &st.Location, &st.Quantity, &st.DateAdded, &st.UpdatedAt); err != nil {
			return nil, err
		}
		stocks = append(stocks, st)
	}
	return stocks, rows.Err()
}

func (s *Store) GetStock(id int) (*model.Stock, error) {
	var st model.Stock
	err := s.db.QueryRow(
		"SELECT id, item_id, location, quantity, date_added, updated_at FROM stock WHERE id = ?", id,
	).Scan(&st.ID, &st.ItemID, &st.Location, &st.Quantity, &st.DateAdded, &st.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func (s *Store) AddStock(st *model.Stock) error {
	res, err := s.db.Exec(
		"INSERT INTO stock (item_id, location, quantity) VALUES (?, ?, ?)",
		st.ItemID, st.Location, st.Quantity)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	st.ID = int(id)
	return nil
}

func (s *Store) UpdateStock(id int, quantity float64) error {
	if quantity <= 0 {
		_, err := s.db.Exec("DELETE FROM stock WHERE id = ?", id)
		return err
	}
	_, err := s.db.Exec("UPDATE stock SET quantity = ?, updated_at = datetime('now') WHERE id = ?", quantity, id)
	return err
}

func (s *Store) DeleteStock(id int) error {
	_, err := s.db.Exec("DELETE FROM stock WHERE id = ?", id)
	return err
}

func (s *Store) MoveStock(fromID int, req model.StockMoveRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromStock model.Stock
	err = tx.QueryRow(
		"SELECT id, item_id, location, quantity FROM stock WHERE id = ?", fromID,
	).Scan(&fromStock.ID, &fromStock.ItemID, &fromStock.Location, &fromStock.Quantity)
	if err != nil {
		return err
	}

	if req.Quantity > fromStock.Quantity {
		return fmt.Errorf("cannot move %.2f, only %.2f available", req.Quantity, fromStock.Quantity)
	}
	if req.ToLocation == fromStock.Location {
		return fmt.Errorf("source and destination are the same")
	}

	newQty := fromStock.Quantity - req.Quantity
	if newQty <= 0 {
		_, err = tx.Exec("DELETE FROM stock WHERE id = ?", fromID)
	} else {
		_, err = tx.Exec("UPDATE stock SET quantity = ?, updated_at = datetime('now') WHERE id = ?", newQty, fromID)
	}
	if err != nil {
		return err
	}

	var existingID int
	err = tx.QueryRow(
		"SELECT id FROM stock WHERE item_id = ? AND location = ?", fromStock.ItemID, req.ToLocation,
	).Scan(&existingID)
	if err == nil {
		_, err = tx.Exec("UPDATE stock SET quantity = quantity + ?, updated_at = datetime('now') WHERE id = ?",
			req.Quantity, existingID)
	} else {
		_, err = tx.Exec("INSERT INTO stock (item_id, location, quantity) VALUES (?, ?, ?)",
			fromStock.ItemID, req.ToLocation, req.Quantity)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}
