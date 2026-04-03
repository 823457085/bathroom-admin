package handler

import (
	"net/http"
	"strconv"

	"bathroom-admin/internal/model"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productRepo *model.ProductRepository
}

func NewProductHandler(pr *model.ProductRepository) *ProductHandler {
	return &ProductHandler{productRepo: pr}
}

type ListProductsRequest struct {
	Page       int    `form:"page" default:"1"`
	PageSize   int    `form:"page_size" default:"20"`
	CategoryID int64  `form:"category_id"`
	Keyword    string `form:"keyword"`
}

type CreateProductRequest struct {
	CategoryID  int64           `json:"category_id" binding:"required"`
	Name        string          `json:"name" binding:"required"`
	Subtitle    string          `json:"subtitle"`
	Price       float64         `json:"price" binding:"required,gt=0"`
	Stock       int             `json:"stock"`
	MainImage   string          `json:"main_image"`
	Images      string          `json:"images"`
	Specs       string          `json:"specs"`
	Description string          `json:"description"`
	Status      int             `json:"status"`
}

func (h *ProductHandler) List(c *gin.Context) {
	var req ListProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	products, total, err := h.productRepo.FindAll(req.Page, req.PageSize, req.CategoryID, req.Keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     req.Page,
		"page_size": req.PageSize,
	})
}

func (h *ProductHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.productRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := &model.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Subtitle:    req.Subtitle,
		Price:       req.Price,
		Stock:       req.Stock,
		MainImage:   req.MainImage,
		Description: req.Description,
		Status:      1,
	}

	id, err := h.productRepo.Create(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := &model.Product{ID: id}
	if req.Name != "" {
		p.Name = req.Name
	}
	p.Subtitle = req.Subtitle
	p.Price = req.Price
	p.Stock = req.Stock
	p.MainImage = req.MainImage
	p.Description = req.Description
	if req.Status > 0 {
		p.Status = req.Status
	}

	if err := h.productRepo.Update(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	if err := h.productRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
