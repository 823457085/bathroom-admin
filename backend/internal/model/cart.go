package model

import (
	"database/sql"
	"time"
)

type CartItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartItemWithProduct struct {
	CartItem
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	MainImage   string  `json:"main_image"`
}

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) FindByUserID(userID int64) ([]CartItemWithProduct, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.user_id, c.product_id, c.quantity, c.created_at, c.updated_at,
		p.name, p.price, p.main_image
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItemWithProduct
	for rows.Next() {
		var ci CartItemWithProduct
		if err := rows.Scan(&ci.ID, &ci.UserID, &ci.ProductID, &ci.Quantity, &ci.CreatedAt, &ci.UpdatedAt, &ci.ProductName, &ci.ProductPrice, &ci.MainImage); err != nil {
			return nil, err
		}
		items = append(items, ci)
	}
	return items, nil
}

func (r *CartRepository) FindByUserAndProduct(userID, productID int64) (*CartItem, error) {
	row := r.db.QueryRow(
		"SELECT id, user_id, product_id, quantity, created_at, updated_at FROM cart_items WHERE user_id = ? AND product_id = ?",
		userID, productID,
	)
	var ci CartItem
	err := row.Scan(&ci.ID, &ci.UserID, &ci.ProductID, &ci.Quantity, &ci.CreatedAt, &ci.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ci, nil
}

func (r *CartRepository) Upsert(ci *CartItem) error {
	if ci.ID > 0 {
		_, err := r.db.Exec("UPDATE cart_items SET quantity = quantity + ? WHERE id = ?", ci.Quantity, ci.ID)
		return err
	}
	_, err := r.db.Exec(
		"INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?, ?, ?)",
		ci.UserID, ci.ProductID, ci.Quantity,
	)
	return err
}

func (r *CartRepository) UpdateQuantity(id int64, quantity int) error {
	_, err := r.db.Exec("UPDATE cart_items SET quantity = ? WHERE id = ?", quantity, id)
	return err
}

func (r *CartRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM cart_items WHERE id = ?", id)
	return err
}

func (r *CartRepository) ClearByUserID(userID int64) error {
	_, err := r.db.Exec("DELETE FROM cart_items WHERE user_id = ?", userID)
	return err
}
