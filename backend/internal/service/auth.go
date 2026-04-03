package service

import (
	"database/sql"
	"errors"
	"bathroom-admin/internal/model"
	"bathroom-admin/internal/repository"
	"bathroom-admin/pkg/jwt"
	"bathroom-admin/pkg/password"
)

var (
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidCreds     = errors.New("invalid phone or password")
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtObj   *jwt.JWT
}

func NewAuthService(userRepo *repository.UserRepository, jwtObj *jwt.JWT) *AuthService {
	return &AuthService{userRepo: userRepo, jwtObj: jwtObj}
}

func (s *AuthService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if user exists
	_, err := s.userRepo.FindByPhone(req.Phone)
	if err == nil {
		return nil, ErrUserExists
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	// Hash password
	pwdHash, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		Phone:        req.Phone,
		PasswordHash: pwdHash,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.jwtObj.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.FindByPhone(req.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCreds
		}
		return nil, err
	}

	if !password.Verify(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCreds
	}

	token, err := s.jwtObj.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) Logout(token string) error {
	// JWT is stateless; logout is handled client-side by discarding the token.
	// For future enhancement: maintain a token blacklist in Redis.
	return nil
}

func (s *AuthService) GetUserByID(id uint64) (*model.User, error) {
	return s.userRepo.FindByID(id)
}
