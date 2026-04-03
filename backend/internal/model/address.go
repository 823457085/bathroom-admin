package model

import (
	"database/sql"
	"time"
)

type Address struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	ReceiverName string    `json:"receiver_name"`
	Phone        string    `json:"phone"`
	Province     string    `json:"province"`
	City         string    `json:"city"`
	District     string    `json:"district"`
	Detail       string    `json:"detail"`
	IsDefault    int       `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AddressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) FindByUserID(userID int64) ([]Address, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, receiver_name, phone, province, city, district, detail, is_default, created_at, updated_at FROM addresses WHERE user_id = ? ORDER BY is_default DESC, id DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []Address
	for rows.Next() {
		var a Address
		if err := rows.Scan(&a.ID, &a.UserID, &a.ReceiverName, &a.Phone, &a.Province, &a.City, &a.District, &a.Detail, &a.IsDefault, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

func (r *AddressRepository) FindByID(id int64) (*Address, error) {
	row := r.db.QueryRow(
		"SELECT id, user_id, receiver_name, phone, province, city, district, detail, is_default, created_at, updated_at FROM addresses WHERE id = ?",
		id,
	)
	var a Address
	err := row.Scan(&a.ID, &a.UserID, &a.ReceiverName, &a.Phone, &a.Province, &a.City, &a.District, &a.Detail, &a.IsDefault, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AddressRepository) Create(a *Address) (int64, error) {
	// 如果设为默认，先清除其他默认
	if a.IsDefault == 1 {
		r.db.Exec("UPDATE addresses SET is_default = 0 WHERE user_id = ?", a.UserID)
	}
	result, err := r.db.Exec(
		"INSERT INTO addresses (user_id, receiver_name, phone, province, city, district, detail, is_default) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		a.UserID, a.ReceiverName, a.Phone, a.Province, a.City, a.District, a.Detail, a.IsDefault,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *AddressRepository) SetDefault(userID, addressID int64) error {
	_, err := r.db.Exec("UPDATE addresses SET is_default = 0 WHERE user_id = ?", userID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("UPDATE addresses SET is_default = 1 WHERE id = ? AND user_id = ?", addressID, userID)
	return err
}
