package handler

import (
	"net/http"

	"bathroom-admin/internal/model"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberRepo *model.MemberRepository
}

func NewMemberHandler(mr *model.MemberRepository) *MemberHandler {
	return &MemberHandler{memberRepo: mr}
}

func (h *MemberHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	member, err := h.memberRepo.GetOrCreate(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get member"})
		return
	}
	levels, _ := h.memberRepo.FindLevels()
	coupons, _ := h.memberRepo.GetCoupons(userID)
	c.JSON(http.StatusOK, gin.H{"member": member, "levels": levels, "coupons": coupons})
}

func (h *MemberHandler) GetLevels(c *gin.Context) {
	levels, err := h.memberRepo.FindLevels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get levels"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"levels": levels})
}

func (h *MemberHandler) ClaimCoupon(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req struct {
		CouponID int64 `json:"coupon_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.memberRepo.ClaimCoupon(userID, req.CouponID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to claim coupon"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "claimed"})
}
