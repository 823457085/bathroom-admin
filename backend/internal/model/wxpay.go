package model

import (
	"database/sql"
	"time"
)

type PayCallback struct {
	ID           int64     `json:"id"`
	OrderID      int64     `json:"order_id"`
	TransactionID string   `json:"transaction_id"`
	Amount       int       `json:"amount"`
	Status       string    `json:"status"`
	RawXML       string    `json:"raw_xml"`
	CreatedAt    time.Time `json:"created_at"`
}

type PayCallbackRepository struct {
	db *sql.DB
}

func NewPayCallbackRepository(db *sql.DB) *PayCallbackRepository {
	return &PayCallbackRepository{db: db}
}

func (r *PayCallbackRepository) Create(cb *PayCallback) error {
	_, err := r.db.Exec(
		"INSERT INTO pay_callbacks (order_id, transaction_id, amount, status, raw_xml) VALUES (?, ?, ?, ?, ?)",
		cb.OrderID, cb.TransactionID, cb.Amount, cb.Status, cb.RawXML,
	)
	return err
}

func (r *PayCallbackRepository) FindByOrderID(orderID int64) (*PayCallback, error) {
	row := r.db.QueryRow(
		"SELECT id, order_id, transaction_id, amount, status, raw_xml, created_at FROM pay_callbacks WHERE order_id = ? ORDER BY id DESC LIMIT 1",
		orderID,
	)
	var cb PayCallback
	err := row.Scan(&cb.ID, &cb.OrderID, &cb.TransactionID, &cb.Amount, &cb.Status, &cb.RawXML, &cb.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &cb, nil
}
