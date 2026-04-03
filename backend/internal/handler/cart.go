package handler

import (
	"net/http"
	"strconv"

	"bathroom-admin/internal/model"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartRepo    *model.CartRepository
	productRepo *model.ProductRepository
}

func NewCartHandler(cr *model.CartRepository, pr *model.ProductRepository) *CartHandler {
	return &CartHandler{cartRepo: cr, productRepo: pr}
}

type AddCartRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

func (h *CartHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	items, err := h.cartRepo.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *CartHandler) Add(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req AddCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查商品是否存在
	_, err := h.productRepo.FindByID(req.ProductID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// 检查是否已在购物车
	existing, err := h.cartRepo.FindByUserAndProduct(userID, req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check cart"})
		return
	}

	ci := &model.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}
	if existing != nil {
		ci.ID = existing.ID
	}

	if err := h.cartRepo.Upsert(ci); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to cart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "added to cart"})
}

func (h *CartHandler) UpdateQuantity(c *gin.Context) {
	idStr := c.Param("item_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartRepo.UpdateQuantity(id, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update quantity"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *CartHandler) Remove(c *gin.Context) {
	idStr := c.Param("item_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	if err := h.cartRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}
