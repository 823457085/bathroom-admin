package repository

import (
	"database/sql"
	"bathroom-admin/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (phone, password_hash, created_at, updated_at) VALUES (?, ?, NOW(), NOW())`
	result, err := r.db.Exec(query, user.Phone, user.PasswordHash)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint64(id)
	return nil
}

func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	query := `SELECT id, phone, email, nickname, avatar, password_hash, created_at, updated_at FROM users WHERE phone = ?`
	row := r.db.QueryRow(query, phone)
	var user model.User
	var email, nickname, avatar sql.NullString
	err := row.Scan(&user.ID, &user.Phone, &email, &nickname, &avatar, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	user.Email = email.String
	user.Nickname = nickname.String
	user.Avatar = avatar.String
	return &user, nil
}

func (r *UserRepository) FindByID(id uint64) (*model.User, error) {
	query := `SELECT id, phone, email, nickname, avatar, password_hash, created_at, updated_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)
	var user model.User
	var email, nickname, avatar sql.NullString
	err := row.Scan(&user.ID, &user.Phone, &email, &nickname, &avatar, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	user.Email = email.String
	user.Nickname = nickname.String
	user.Avatar = avatar.String
	return &user, nil
}
