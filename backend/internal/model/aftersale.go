package model

import (
	"database/sql"
	"time"
)

type Aftersale struct {
	ID          int64     `json:"id"`
	OrderID    int64     `json:"order_id"`
	UserID     int64     `json:"user_id"`
	Type       int       `json:"type"`        // 1: 退货退款, 2: 仅退款
	Status     int       `json:"status"`      // 1: 待处理, 2: 已同意, 3: 已拒绝, 4: 已完成
	Reason     string    `json:"reason"`
	Amount     float64   `json:"amount"`
	Description string   `json:"description"`
	Images     string    `json:"images"`      // JSON array of image URLs
	Reply      string    `json:"reply"`       // 商家回复
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AftersaleRepository struct {
	db *sql.DB
}

func NewAftersaleRepository(db *sql.DB) *AftersaleRepository {
	return &AftersaleRepository{db: db}
}

func (r *AftersaleRepository) Create(a *Aftersale) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO aftersales (order_id, user_id, type, status, reason, amount, description, images) VALUES (?, ?, ?, 1, ?, ?, ?, ?)",
		a.OrderID, a.UserID, a.Type, a.Reason, a.Amount, a.Description, a.Images,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *AftersaleRepository) FindByUserID(userID int64) ([]Aftersale, error) {
	rows, err := r.db.Query(
		"SELECT id, order_id, user_id, type, status, reason, amount, description, images, reply, created_at, updated_at FROM aftersales WHERE user_id = ? ORDER BY id DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Aftersale
	for rows.Next() {
		var a Aftersale
		if err := rows.Scan(&a.ID, &a.OrderID, &a.UserID, &a.Type, &a.Status, &a.Reason, &a.Amount, &a.Description, &a.Images, &a.Reply, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *AftersaleRepository) FindAll(page, pageSize int) ([]Aftersale, int, error) {
	var total int
	r.db.QueryRow("SELECT COUNT(*) FROM aftersales").Scan(&total)
	rows, err := r.db.Query(
		"SELECT id, order_id, user_id, type, status, reason, amount, description, images, reply, created_at, updated_at FROM aftersales ORDER BY id DESC LIMIT ? OFFSET ?",
		pageSize, (page-1)*pageSize,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []Aftersale
	for rows.Next() {
		var a Aftersale
		if err := rows.Scan(&a.ID, &a.OrderID, &a.UserID, &a.Type, &a.Status, &a.Reason, &a.Amount, &a.Description, &a.Images, &a.Reply, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, a)
	}
	return list, total, nil
}

func (r *AftersaleRepository) UpdateStatus(id int64, status int, reply string) error {
	_, err := r.db.Exec("UPDATE aftersales SET status = ?, reply = ? WHERE id = ?", status, reply, id)
	return err
}
