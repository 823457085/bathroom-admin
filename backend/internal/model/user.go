package model

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (phone, password_hash, created_at, updated_at) VALUES (?, ?, NOW(), NOW())",
		user.Phone, user.PasswordHash,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepository) FindByPhone(phone string) (*User, error) {
	row := r.db.QueryRow(
		"SELECT id, phone, email, nickname, avatar, password_hash, created_at, updated_at FROM users WHERE phone = ?",
		phone,
	)
	var u User
	err := row.Scan(&u.ID, &u.Phone, &u.Email, &u.Nickname, &u.Avatar, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id int64) (*User, error) {
	row := r.db.QueryRow(
		"SELECT id, phone, email, nickname, avatar, password_hash, created_at, updated_at FROM users WHERE id = ?",
		id,
	)
	var u User
	err := row.Scan(&u.ID, &u.Phone, &u.Email, &u.Nickname, &u.Avatar, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
