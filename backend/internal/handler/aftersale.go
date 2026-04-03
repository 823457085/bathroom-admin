package handler

import (
	"net/http"
	"strconv"

	"bathroom-admin/internal/model"

	"github.com/gin-gonic/gin"
)

type AftersaleHandler struct {
	aftersaleRepo *model.AftersaleRepository
}

func NewAftersaleHandler(ar *model.AftersaleRepository) *AftersaleHandler {
	return &AftersaleHandler{aftersaleRepo: ar}
}

type CreateAftersaleRequest struct {
	OrderID     int64   `json:"order_id" binding:"required"`
	Type       int     `json:"type" binding:"required"` // 1退货退款 2仅退款
	Reason     string  `json:"reason" binding:"required"`
	Amount     float64 `json:"amount"`
	Description string `json:"description"`
	Images     string  `json:"images"`
}

func (h *AftersaleHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req CreateAftersaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a := &model.Aftersale{
		OrderID:     req.OrderID,
		UserID:     userID,
		Type:       req.Type,
		Reason:     req.Reason,
		Amount:     req.Amount,
		Description: req.Description,
		Images:     req.Images,
	}
	id, err := h.aftersaleRepo.Create(a)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create aftersale"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *AftersaleHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	list, err := h.aftersaleRepo.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get list"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list})
}

// Admin: list all aftersales
func (h *AftersaleHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	list, total, err := h.aftersaleRepo.FindAll(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get list"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total, "page": page})
}

// Admin: handle aftersale
func (h *AftersaleHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Status int    `json:"status" binding:"required"` // 2同意 3拒绝
		Reply  string `json:"reply"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.aftersaleRepo.UpdateStatus(id, req.Status, req.Reply); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to handle"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
