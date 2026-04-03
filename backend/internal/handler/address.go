package handler

import (
	"net/http"
	"strconv"

	"bathroom-admin/internal/model"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addressRepo *model.AddressRepository
}

func NewAddressHandler(ar *model.AddressRepository) *AddressHandler {
	return &AddressHandler{addressRepo: ar}
}

type CreateAddressRequest struct {
	ReceiverName string `json:"receiver_name" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Province     string `json:"province" binding:"required"`
	City         string `json:"city" binding:"required"`
	District     string `json:"district" binding:"required"`
	Detail       string `json:"detail" binding:"required"`
	IsDefault    int    `json:"is_default"`
}

func (h *AddressHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	addresses, err := h.addressRepo.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get addresses"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

func (h *AddressHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addr := &model.Address{
		UserID:       userID,
		ReceiverName: req.ReceiverName,
		Phone:        req.Phone,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Detail:       req.Detail,
		IsDefault:    req.IsDefault,
	}

	id, err := h.addressRepo.Create(addr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create address"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *AddressHandler) SetDefault(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
		return
	}
	if err := h.addressRepo.SetDefault(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set default"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
