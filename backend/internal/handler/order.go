package handler

import (
	"net/http"
	"strconv"

	"bathroom-admin/internal/model"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderRepo   *model.OrderRepository
	cartRepo   *model.CartRepository
	addressRepo *model.AddressRepository
}

func NewOrderHandler(or *model.OrderRepository, cr *model.CartRepository, ar *model.AddressRepository) *OrderHandler {
	return &OrderHandler{orderRepo: or, cartRepo: cr, addressRepo: ar}
}

type CreateOrderRequest struct {
	AddressID int64  `json:"address_id" binding:"required"`
	Remark    string `json:"remark"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取购物车商品
	cartItems, err := h.cartRepo.FindByUserID(userID)
	if err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	// 获取地址
	address, err := h.addressRepo.FindByID(req.AddressID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address not found"})
		return
	}

	// 计算总金额并构建订单项
	var totalAmount float64
	var orderItems []model.OrderItem
	for _, item := range cartItems {
		subtotal := item.ProductPrice * float64(item.Quantity)
		totalAmount += subtotal
		orderItems = append(orderItems, model.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Price:      item.ProductPrice,
			Quantity:   item.Quantity,
			Subtotal:   subtotal,
		})
	}

	order := &model.Order{
		OrderNo:     model.GenerateOrderNo(),
		UserID:      userID,
		AddressID:   address.ID,
		TotalAmount: totalAmount,
		Status:      1, // 待付款
		Remark:      req.Remark,
	}

	orderID, err := h.orderRepo.Create(order, orderItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	// 清空购物车
	h.cartRepo.ClearByUserID(userID)

	c.JSON(http.StatusOK, gin.H{"order_id": orderID, "order_no": order.OrderNo})
}

func (h *OrderHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := h.orderRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders, "total": total, "page": page, "page_size": pageSize})
}

func (h *OrderHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.orderRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	items, err := h.orderRepo.FindItemsByOrderID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order items"})
		return
	}

	address, _ := h.addressRepo.FindByID(order.AddressID)

	c.JSON(http.StatusOK, gin.H{"order": order, "items": items, "address": address})
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.orderRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if order.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only pending orders can be cancelled"})
		return
	}

	if err := h.orderRepo.UpdateStatus(id, 2); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}
