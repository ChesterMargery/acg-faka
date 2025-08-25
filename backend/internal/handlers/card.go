package handlers

import (
	"acg-faka/backend/internal/database"
	"acg-faka/backend/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CardHandler struct{}

type CreateCardRequest struct {
	CommodityID uint   `json:"commodity_id" binding:"required"`
	Draft       string `json:"draft" binding:"required"`
	Secret      string `json:"secret"`
	Note        string `json:"note"`
	Race        string `json:"race"`
}

type BatchCreateCardRequest struct {
	CommodityID uint     `json:"commodity_id" binding:"required"`
	Cards       []string `json:"cards" binding:"required"`
	Race        string   `json:"race"`
}

type UpdateCardRequest struct {
	Draft  string `json:"draft"`
	Secret string `json:"secret"`
	Note   string `json:"note"`
	Race   string `json:"race"`
	Status int    `json:"status"`
}

// GetCards returns all cards with pagination
func (h *CardHandler) GetCards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var cards []models.Card
	var total int64

	query := database.DB.Model(&models.Card{}).Preload("OwnerUser").Preload("Commodity").Preload("Order")

	// Add filters
	if commodityID := c.Query("commodity_id"); commodityID != "" {
		query = query.Where("commodity_id = ?", commodityID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if owner := c.Query("owner"); owner != "" {
		query = query.Where("owner = ?", owner)
	}
	if race := c.Query("race"); race != "" {
		query = query.Where("race = ?", race)
	}

	// Count total
	query.Count(&total)

	// Get data with pagination
	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"data":  cards,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetCard returns a specific card
func (h *CardHandler) GetCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	var card models.Card
	if err := database.DB.Preload("OwnerUser").Preload("Commodity").Preload("Order").First(&card, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": card,
	})
}

// CreateCard creates a new card
func (h *CardHandler) CreateCard(c *gin.Context) {
	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	// Check if commodity exists and user has permission
	var commodity models.Commodity
	if err := database.DB.First(&commodity, req.CommodityID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Commodity not found"})
		return
	}

	// Check ownership
	if userType != "admin" && commodity.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	card := models.Card{
		CommodityID: req.CommodityID,
		Draft:       req.Draft,
		Secret:      req.Secret,
		Note:        req.Note,
		Race:        req.Race,
		Owner:       userID,
		CreateTime:  time.Now(),
		Status:      0, // 0: available
	}

	if err := database.DB.Create(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create card"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    200,
		"data":    card,
		"message": "Card created successfully",
	})
}

// BatchCreateCards creates multiple cards at once
func (h *CardHandler) BatchCreateCards(c *gin.Context) {
	var req BatchCreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	// Check if commodity exists and user has permission
	var commodity models.Commodity
	if err := database.DB.First(&commodity, req.CommodityID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Commodity not found"})
		return
	}

	// Check ownership
	if userType != "admin" && commodity.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var cards []models.Card
	now := time.Now()

	for _, cardData := range req.Cards {
		card := models.Card{
			CommodityID: req.CommodityID,
			Draft:       cardData,
			Race:        req.Race,
			Owner:       userID,
			CreateTime:  now,
			Status:      0, // 0: available
		}
		cards = append(cards, card)
	}

	if err := database.DB.CreateInBatches(cards, 100).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cards"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    200,
		"data":    cards,
		"message": "Cards created successfully",
		"count":   len(cards),
	})
}

// UpdateCard updates an existing card
func (h *CardHandler) UpdateCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	var req UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	var card models.Card
	if err := database.DB.First(&card, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	// Check ownership
	if userType != "admin" && card.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update fields
	if req.Draft != "" {
		card.Draft = req.Draft
	}
	card.Secret = req.Secret
	card.Note = req.Note
	card.Race = req.Race
	if req.Status >= 0 {
		card.Status = req.Status
	}

	if err := database.DB.Save(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    card,
		"message": "Card updated successfully",
	})
}

// DeleteCard deletes a card
func (h *CardHandler) DeleteCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	var card models.Card
	if err := database.DB.First(&card, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	// Check ownership
	if userType != "admin" && card.Owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if card is sold
	if card.Status == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete sold card"})
		return
	}

	if err := database.DB.Delete(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Card deleted successfully",
	})
}

// BatchDeleteCards deletes multiple cards
func (h *CardHandler) BatchDeleteCards(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	// Check ownership and status
	var cards []models.Card
	query := database.DB.Where("id IN ?", req.IDs)
	if userType != "admin" {
		query = query.Where("owner = ?", userID)
	}
	
	if err := query.Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if len(cards) != len(req.IDs) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Some cards not found or access denied"})
		return
	}

	// Check if any card is sold
	for _, card := range cards {
		if card.Status == 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete sold cards"})
			return
		}
	}

	if err := database.DB.Delete(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Cards deleted successfully",
		"count":   len(cards),
	})
}