package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Order struct {
	ID          int64     `json:"id"`
	OrderNo     string    `json:"order_no"`
	UserID      int64     `json:"user_id"`
	AddressID   int64     `json:"address_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      int       `json:"status"`
	Remark      string    `json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID          int64     `json:"id"`
	OrderID     int64     `json:"order_id"`
	ProductID   int64     `json:"product_id"`
	ProductName string    `json:"product_name"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Subtotal    float64   `json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(o *Order, items []OrderItem) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.Exec(
		"INSERT INTO orders (order_no, user_id, address_id, total_amount, status, remark) VALUES (?, ?, ?, ?, ?, ?)",
		o.OrderNo, o.UserID, o.AddressID, o.TotalAmount, o.Status, o.Remark,
	)
	if err != nil {
		return 0, err
	}
	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		_, err := tx.Exec(
			"INSERT INTO order_items (order_id, product_id, product_name, price, quantity, subtotal) VALUES (?, ?, ?, ?, ?, ?)",
			orderID, item.ProductID, item.ProductName, item.Price, item.Quantity, item.Subtotal,
		)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	return orderID, err
}

func (r *OrderRepository) FindByUserID(userID int64, page, pageSize int) ([]Order, int, error) {
	var total int
	countRow := r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE user_id = ?", userID)
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		"SELECT id, order_no, user_id, address_id, total_amount, status, remark, created_at, updated_at FROM orders WHERE user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?",
		userID, pageSize, (page-1)*pageSize,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.OrderNo, &o.UserID, &o.AddressID, &o.TotalAmount, &o.Status, &o.Remark, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}
	return orders, total, nil
}

func (r *OrderRepository) FindByID(id int64) (*Order, error) {
	row := r.db.QueryRow(
		"SELECT id, order_no, user_id, address_id, total_amount, status, remark, created_at, updated_at FROM orders WHERE id = ?",
		id,
	)
	var o Order
	err := row.Scan(&o.ID, &o.OrderNo, &o.UserID, &o.AddressID, &o.TotalAmount, &o.Status, &o.Remark, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) FindByOrderNo(orderNo string) (*Order, error) {
	row := r.db.QueryRow(
		"SELECT id, order_no, user_id, address_id, total_amount, status, remark, created_at, updated_at FROM orders WHERE order_no = ?",
		orderNo,
	)
	var o Order
	err := row.Scan(&o.ID, &o.OrderNo, &o.UserID, &o.AddressID, &o.TotalAmount, &o.Status, &o.Remark, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) FindItemsByOrderID(orderID int64) ([]OrderItem, error) {
	rows, err := r.db.Query(
		"SELECT id, order_id, product_id, product_name, price, quantity, subtotal, created_at FROM order_items WHERE order_id = ?",
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.Price, &item.Quantity, &item.Subtotal, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OrderRepository) UpdateStatus(id int64, status int) error {
	_, err := r.db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, id)
	return err
}

func GenerateOrderNo() string {
	return fmt.Sprintf("%d%06d", time.Now().UnixNano()/1000000, time.Now().Second()*1000+int(time.Now().Nanosecond()/1000000))
}
