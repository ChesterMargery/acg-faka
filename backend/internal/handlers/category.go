package handlers

import (
	"acg-faka/backend/internal/database"
	"acg-faka/backend/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct{}

type CreateCategoryRequest struct {
	Name   string `json:"name" binding:"required"`
	Icon   string `json:"icon"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

type UpdateCategoryRequest struct {
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

// GetCategories returns all categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	query := database.DB.Preload("OwnerUser").Order("sort ASC, id DESC")

	// Add filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": categories,
	})
}

// GetCategory returns a specific category
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category models.Category
	if err := database.DB.Preload("OwnerUser").Preload("Commodities").First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": category,
	})
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	// Only admin or business users can create categories
	if userType != "admin" {
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if user.BusinessLevel == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Business level required"})
			return
		}
	}

	category := models.Category{
		Name:       req.Name,
		Icon:       req.Icon,
		Sort:       req.Sort,
		Owner:      userID,
		Status:     req.Status,
		CreateTime: time.Now(),
	}

	if req.Status == 0 {
		category.Status = 1 // Default to active
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 200,
		"data": category,
		"message": "Category created successfully",
	})
}

// UpdateCategory updates an existing category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check ownership
	if userType != "admin" && category.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	category.Sort = req.Sort
	if req.Status != 0 {
		category.Status = req.Status
	}

	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": category,
		"message": "Category updated successfully",
	})
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check ownership
	if userType != "admin" && category.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if category has commodities
	var commodityCount int64
	database.DB.Model(&models.Commodity{}).Where("category_id = ?", id).Count(&commodityCount)
	if commodityCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with commodities"})
		return
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Category deleted successfully",
	})
}