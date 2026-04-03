package model

import (
	"database/sql"
	"time"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  int64     `json:"parent_id"`
	Sort      int       `json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]Category, error) {
	rows, err := r.db.Query("SELECT id, name, parent_id, sort, created_at, updated_at FROM categories ORDER BY sort ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID, &c.Sort, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) FindByID(id int64) (*Category, error) {
	row := r.db.QueryRow("SELECT id, name, parent_id, sort, created_at, updated_at FROM categories WHERE id = ?", id)
	var c Category
	err := row.Scan(&c.ID, &c.Name, &c.ParentID, &c.Sort, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(c *Category) (int64, error) {
	result, err := r.db.Exec("INSERT INTO categories (name, parent_id, sort) VALUES (?, ?, ?)", c.Name, c.ParentID, c.Sort)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
