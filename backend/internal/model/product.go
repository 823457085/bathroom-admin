package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Product struct {
	ID          int64           `json:"id"`
	CategoryID  int64           `json:"category_id"`
	Name        string          `json:"name"`
	Subtitle    string          `json:"subtitle"`
	Price       float64         `json:"price"`
	Stock       int             `json:"stock"`
	MainImage   string          `json:"main_image"`
	Images      json.RawMessage `json:"images"`
	Specs       json.RawMessage `json:"specs"`
	Description string          `json:"description"`
	Status      int             `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll(page, pageSize int, categoryID int64, keyword string) ([]Product, int, error) {
	countSQL := "SELECT COUNT(*) FROM products WHERE 1=1"
	listSQL := "SELECT id, category_id, name, subtitle, price, stock, main_image, images, specs, description, status, created_at, updated_at FROM products WHERE 1=1"
	var args []interface{}

	if categoryID > 0 {
		countSQL += " AND category_id = ?"
		listSQL += " AND category_id = ?"
		args = append(args, categoryID)
	}
	if keyword != "" {
		kw := "%" + keyword + "%"
		countSQL += " AND name LIKE ?"
		listSQL += " AND name LIKE ?"
		args = append(args, kw)
	}

	// total count
	var total int
	countRow := r.db.QueryRow(countSQL, args...)
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, err
	}

	// list with pagination
	listSQL += " AND status = 1 ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Subtitle, &p.Price, &p.Stock, &p.MainImage, &p.Images, &p.Specs, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	return products, total, nil
}

func (r *ProductRepository) FindByID(id int64) (*Product, error) {
	row := r.db.QueryRow(
		"SELECT id, category_id, name, subtitle, price, stock, main_image, images, specs, description, status, created_at, updated_at FROM products WHERE id = ?",
		id,
	)
	var p Product
	err := row.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Subtitle, &p.Price, &p.Stock, &p.MainImage, &p.Images, &p.Specs, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(p *Product) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO products (category_id, name, subtitle, price, stock, main_image, images, specs, description, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		p.CategoryID, p.Name, p.Subtitle, p.Price, p.Stock, p.MainImage, p.Images, p.Specs, p.Description, p.Status,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ProductRepository) Update(p *Product) error {
	fields := []string{}
	args := []interface{}{}

	if p.Name != "" {
		fields = append(fields, "name = ?")
		args = append(args, p.Name)
	}
	if p.Subtitle != "" {
		fields = append(fields, "subtitle = ?")
		args = append(args, p.Subtitle)
	}
	if p.Price > 0 {
		fields = append(fields, "price = ?")
		args = append(args, p.Price)
	}
	if p.Stock >= 0 {
		fields = append(fields, "stock = ?")
		args = append(args, p.Stock)
	}
	if p.MainImage != "" {
		fields = append(fields, "main_image = ?")
		args = append(args, p.MainImage)
	}
	if p.Description != "" {
		fields = append(fields, "description = ?")
		args = append(args, p.Description)
	}
	if p.Status > 0 {
		fields = append(fields, "status = ?")
		args = append(args, p.Status)
	}

	if len(fields) == 0 {
		return nil
	}

	args = append(args, p.ID)
	_, err := r.db.Exec("UPDATE products SET "+strings.Join(fields, ", ")+" WHERE id = ?", args...)
	return err
}

func (r *ProductRepository) Delete(id int64) error {
	_, err := r.db.Exec("UPDATE products SET status = 0 WHERE id = ?", id)
	return err
}

func (r *ProductRepository) IsAdmin() bool { return true }

func (r *ProductRepository) FindAllAdmin(categoryID int64, keyword string) ([]Product, int, error) {
	return r.FindAll(1, 1000, categoryID, keyword)
}
