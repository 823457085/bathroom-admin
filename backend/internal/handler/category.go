package handler

import (
	"net/http"

	"bathroom-admin/internal/model"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryRepo *model.CategoryRepository
}

func NewCategoryHandler(cr *model.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{categoryRepo: cr}
}

type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID int64  `json:"parent_id"`
	Sort     int    `json:"sort"`
}

func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get categories"})
		return
	}
	// Build tree
	tree := buildCategoryTree(categories, 0)
	c.JSON(http.StatusOK, gin.H{"categories": tree})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat := &model.Category{Name: req.Name, ParentID: req.ParentID, Sort: req.Sort}
	id, err := h.categoryRepo.Create(cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

type treeNode struct {
	model.Category
	Children []treeNode `json:"children,omitempty"`
}

func buildCategoryTree(categories []model.Category, parentID int64) []treeNode {
	var result []treeNode
	for _, c := range categories {
		if c.ParentID == parentID {
			node := treeNode{Category: c}
			node.Children = buildCategoryTree(categories, c.ID)
			result = append(result, node)
		}
	}
	return result
}
