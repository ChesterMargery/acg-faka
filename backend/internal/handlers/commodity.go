package handlers

import (
	"acg-faka/backend/internal/database"
	"acg-faka/backend/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommodityHandler struct{}

type CreateCommodityRequest struct {
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Cover          string  `json:"cover"`
	Price          float64 `json:"price" binding:"required,gt=0"`
	VipPrice       float64 `json:"vip_price"`
	CategoryID     uint    `json:"category_id" binding:"required"`
	Status         int     `json:"status"`
	Sort           int     `json:"sort"`
	IsAutoDelivery bool    `json:"is_auto_delivery"`
}

type UpdateCommodityRequest struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Cover          string  `json:"cover"`
	Price          float64 `json:"price"`
	VipPrice       float64 `json:"vip_price"`
	CategoryID     uint    `json:"category_id"`
	Status         int     `json:"status"`
	Sort           int     `json:"sort"`
	IsAutoDelivery bool    `json:"is_auto_delivery"`
}

// GetCommodities returns all commodities with pagination
func (h *CommodityHandler) GetCommodities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var commodities []models.Commodity
	var total int64

	query := database.DB.Model(&models.Commodity{}).Preload("Category").Preload("OwnerUser")

	// Add filters
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if owner := c.Query("owner"); owner != "" {
		query = query.Where("owner = ?", owner)
	}
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Count total
	query.Count(&total)

	// Get data with pagination
	if err := query.Order("sort ASC, id DESC").Offset(offset).Limit(limit).Find(&commodities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get card counts for each commodity
	for i := range commodities {
		var cardCount int64
		database.DB.Model(&models.Card{}).Where("commodity_id = ? AND status = 0", commodities[i].ID).Count(&cardCount)
		// Add card count as a custom field (using a map to avoid modifying struct)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"data":  commodities,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetCommodity returns a specific commodity
func (h *CommodityHandler) GetCommodity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid commodity ID"})
		return
	}

	var commodity models.Commodity
	if err := database.DB.Preload("Category").Preload("OwnerUser").First(&commodity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commodity not found"})
		return
	}

	// Get available card count
	var cardCount int64
	database.DB.Model(&models.Card{}).Where("commodity_id = ? AND status = 0", id).Count(&cardCount)

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"data":      commodity,
		"card_count": cardCount,
	})
}

// CreateCommodity creates a new commodity
func (h *CommodityHandler) CreateCommodity(c *gin.Context) {
	var req CreateCommodityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	// Only admin or business users can create commodities
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

	// Check if category exists
	var category models.Category
	if err := database.DB.First(&category, req.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	commodity := models.Commodity{
		Name:           req.Name,
		Description:    req.Description,
		Cover:          req.Cover,
		Price:          req.Price,
		VipPrice:       req.VipPrice,
		CategoryID:     req.CategoryID,
		Owner:          userID,
		Status:         req.Status,
		Sort:           req.Sort,
		IsAutoDelivery: req.IsAutoDelivery,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}

	if req.Status == 0 {
		commodity.Status = 1 // Default to active
	}

	if err := database.DB.Create(&commodity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create commodity"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    200,
		"data":    commodity,
		"message": "Commodity created successfully",
	})
}

// UpdateCommodity updates an existing commodity
func (h *CommodityHandler) UpdateCommodity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid commodity ID"})
		return
	}

	var req UpdateCommodityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	var commodity models.Commodity
	if err := database.DB.First(&commodity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commodity not found"})
		return
	}

	// Check ownership
	if userType != "admin" && commodity.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update fields
	if req.Name != "" {
		commodity.Name = req.Name
	}
	commodity.Description = req.Description
	commodity.Cover = req.Cover
	if req.Price > 0 {
		commodity.Price = req.Price
	}
	commodity.VipPrice = req.VipPrice
	if req.CategoryID > 0 {
		// Check if category exists
		var category models.Category
		if err := database.DB.First(&category, req.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
		commodity.CategoryID = req.CategoryID
	}
	if req.Status > 0 {
		commodity.Status = req.Status
	}
	commodity.Sort = req.Sort
	commodity.IsAutoDelivery = req.IsAutoDelivery
	commodity.UpdateTime = time.Now()

	if err := database.DB.Save(&commodity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update commodity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    commodity,
		"message": "Commodity updated successfully",
	})
}

// DeleteCommodity deletes a commodity
func (h *CommodityHandler) DeleteCommodity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid commodity ID"})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	var commodity models.Commodity
	if err := database.DB.First(&commodity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commodity not found"})
		return
	}

	// Check ownership
	if userType != "admin" && commodity.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if commodity has cards
	var cardCount int64
	database.DB.Model(&models.Card{}).Where("commodity_id = ?", id).Count(&cardCount)
	if cardCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete commodity with cards"})
		return
	}

	if err := database.DB.Delete(&commodity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete commodity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Commodity deleted successfully",
	})
}