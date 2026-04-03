package model

import (
	"database/sql"
	"time"
)

type MemberLevel struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	MinPoints   int       `json:"min_points"`
	Discount    float64   `json:"discount"` // 折扣率，如 0.9 表示 9 折
	CreatedAt   time.Time `json:"created_at"`
}

type Member struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	LevelID   int64     `json:"level_id"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Coupon struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // "cash": 满减券, "discount": 折扣券
	Threshold   float64   `json:"threshold"` // 满多少
	Discount    float64   `json:"discount"` // 减多少 / 折扣率
	TotalCount  int       `json:"total_count"`
	LeftCount   int       `json:"left_count"`
	ExpireAt    time.Time `json:"expire_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserCoupon struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CouponID  int64     `json:"coupon_id"`
	Used      bool      `json:"used"`
	UsedAt    *time.Time `json:"used_at"`
	OrderID   *int64    `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
}

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) FindLevels() ([]MemberLevel, error) {
	rows, err := r.db.Query("SELECT id, name, min_points, discount, created_at FROM member_levels ORDER BY min_points ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var levels []MemberLevel
	for rows.Next() {
		var l MemberLevel
		if err := rows.Scan(&l.ID, &l.Name, &l.MinPoints, &l.Discount, &l.CreatedAt); err != nil {
			return nil, err
		}
		levels = append(levels, l)
	}
	return levels, nil
}

func (r *MemberRepository) GetOrCreate(userID int64) (*Member, error) {
	row := r.db.QueryRow("SELECT id, user_id, level_id, points, created_at, updated_at FROM members WHERE user_id = ?", userID)
	var m Member
	err := row.Scan(&m.ID, &m.UserID, &m.LevelID, &m.Points, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		_, err = r.db.Exec("INSERT INTO members (user_id, level_id, points) VALUES (?, 1, 0)", userID)
		if err != nil {
			return nil, err
		}
		return &Member{UserID: userID, LevelID: 1, Points: 0}, nil
	}
	return &m, err
}

func (r *MemberRepository) AddPoints(userID int64, points int) error {
	_, err := r.db.Exec("UPDATE members SET points = points + ? WHERE user_id = ?", points, userID)
	return err
}

func (r *MemberRepository) GetCoupons(userID int64) ([]UserCoupon, error) {
	rows, err := r.db.Query(
		`SELECT uc.id, uc.user_id, uc.coupon_id, uc.used, uc.used_at, uc.order_id, uc.created_at,
		c.name, c.type, c.threshold, c.discount, c.expire_at
		FROM user_coupons uc
		JOIN coupons c ON uc.coupon_id = c.id
		WHERE uc.user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var coupons []UserCoupon
	for rows.Next() {
		var uc UserCoupon
		if err := rows.Scan(&uc.ID, &uc.UserID, &uc.CouponID, &uc.Used, &uc.UsedAt, &uc.OrderID, &uc.CreatedAt); err != nil {
			return nil, err
		}
		coupons = append(coupons, uc)
	}
	return coupons, nil
}

func (r *MemberRepository) ClaimCoupon(userID, couponID int64) error {
	result, err := r.db.Exec(
		"INSERT INTO user_coupons (user_id, coupon_id) VALUES (?, ?)",
		userID, couponID,
	)
	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	return err
}
