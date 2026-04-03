package handler

import (
	"net/http"
	"sync"
	"time"

	"bathroom-admin/internal/model"
	"bathroom-admin/pkg/jwt"
	"bathroom-admin/pkg/password"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo  *model.UserRepository
	jwtMgr    *jwt.JWTManager
	blacklist sync.Map // token string -> bool
}

func NewAuthHandler(ur *model.UserRepository, jm *jwt.JWTManager) *AuthHandler {
	return &AuthHandler{userRepo: ur, jwtMgr: jm}
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &model.User{Phone: req.Phone, PasswordHash: hash}
	id, err := h.userRepo.Create(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone already registered"})
		return
	}

	token, err := h.jwtMgr.Generate(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": id})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.FindByPhone(req.Phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid phone or password"})
		return
	}

	if !password.Verify(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid phone or password"})
		return
	}

	token, err := h.jwtMgr.Generate(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": user.ID})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := splitAuthHeader(authHeader)
		if len(parts) == 2 {
			h.blacklist.Store(parts[1], time.Now())
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) IsBlacklisted(token string) bool {
	_, ok := h.blacklist.Load(token)
	return ok
}

func splitAuthHeader(authHeader string) []string {
	for i := 0; i < len(authHeader); i++ {
		if authHeader[i] == ' ' {
			return []string{authHeader[:i], authHeader[i+1:]}
		}
	}
	return nil
}
