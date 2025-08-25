package handlers

import (
	"acg-faka/internal/auth"
	"acg-faka/internal/config"
	"acg-faka/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB:  db,
		Cfg: config.Load(),
	}
}

func SetupRoutes(router *gin.RouterGroup, db *gorm.DB) {
	handler := NewHandler(db)

	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/admin/login", handler.AdminLogin)
	}

	// Public API routes
	public := router.Group("/public")
	{
		public.GET("/categories", handler.GetCategories)
		public.GET("/commodities", handler.GetCommodities)
		public.GET("/commodity/:id", handler.GetCommodity)
		public.GET("/configs", handler.GetPublicConfigs)
	}

	// Protected user routes
	user := router.Group("/user")
	user.Use(handler.AuthMiddleware("user"))
	{
		user.GET("/profile", handler.GetProfile)
		user.PUT("/profile", handler.UpdateProfile)
		user.GET("/orders", handler.GetUserOrders)
		user.POST("/orders", handler.CreateOrder)
		user.GET("/balance", handler.GetBalance)
	}

	// Protected admin routes
	admin := router.Group("/admin")
	admin.Use(handler.AuthMiddleware("admin"))
	{
		admin.GET("/dashboard", handler.GetDashboard)
		admin.GET("/users", handler.GetUsers)
		admin.GET("/orders", handler.GetAllOrders)
		admin.GET("/commodities", handler.GetAllCommodities)
		admin.POST("/commodities", handler.CreateCommodity)
		admin.PUT("/commodities/:id", handler.UpdateCommodity)
		admin.DELETE("/commodities/:id", handler.DeleteCommodity)
		admin.GET("/categories", handler.GetAllCategories)
		admin.POST("/categories", handler.CreateCategory)
		admin.PUT("/categories/:id", handler.UpdateCategory)
		admin.DELETE("/categories/:id", handler.DeleteCategory)
	}
}

// AuthMiddleware validates JWT tokens and checks user roles
func (h *Handler) AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := auth.ParseJWT(parts[1], h.Cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if requiredRole != "" && claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// Auth handlers

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := h.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email already exists (if provided)
	if req.Email != nil {
		if err := h.DB.Where("email = ?", *req.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
	}

	// Generate salt and hash password
	salt, err := auth.GenerateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate salt"})
		return
	}

	hashedPassword := auth.HashPassword(req.Password, salt)

	// Generate app key
	appKey, err := auth.GenerateAppKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate app key"})
		return
	}

	// Create user
	user := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Phone:      req.Phone,
		QQ:         req.QQ,
		Password:   hashedPassword,
		Salt:       salt,
		AppKey:     appKey,
		Status:     1, // Active by default
		CreateTime: time.Now(),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Username, "user", h.Cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  &user,
	})
}

// Login handles user login
// @Summary User login
// @Description Login with username/email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	query := h.DB.Where("username = ?", req.Username)
	if strings.Contains(req.Username, "@") {
		query = h.DB.Where("email = ?", req.Username)
	}

	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !auth.VerifyPassword(req.Password, user.Salt, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Status != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is disabled"})
		return
	}

	// Update login info
	now := time.Now()
	user.LastLoginTime = user.LoginTime
	user.LoginTime = &now
	user.LastLoginIP = user.LoginIP
	loginIP := c.ClientIP()
	user.LoginIP = &loginIP

	h.DB.Save(&user)

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Username, "user", h.Cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  &user,
	})
}

// AdminLogin handles admin login
// @Summary Admin login
// @Description Login as admin with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Admin login credentials"
// @Success 200 {object} AdminAuthResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/admin/login [post]
func (h *Handler) AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var manager models.Manager
	if err := h.DB.Where("email = ?", req.Username).First(&manager).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !auth.VerifyPassword(req.Password, manager.Salt, manager.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if manager.Status != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is disabled"})
		return
	}

	// Update login info
	now := time.Now()
	manager.LastLoginTime = manager.LoginTime
	manager.LoginTime = &now
	manager.LastLoginIP = manager.LoginIP
	loginIP := c.ClientIP()
	manager.LoginIP = &loginIP

	h.DB.Save(&manager)

	// Generate JWT token
	token, err := auth.GenerateJWT(manager.ID, manager.Email, "admin", h.Cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AdminAuthResponse{
		Token:   token,
		Manager: &manager,
	})
}

// Public handlers

// GetCategories returns all active categories
// @Summary Get categories
// @Description Get all active product categories
// @Tags public
// @Produce json
// @Success 200 {array} models.Category
// @Router /public/categories [get]
func (h *Handler) GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := h.DB.Where("status = ? AND hide = ?", 1, 0).Order("sort ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCommodities returns active commodities with pagination
// @Summary Get commodities
// @Description Get active commodities with optional category filter and pagination
// @Tags public
// @Produce json
// @Param category_id query int false "Category ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} PaginatedResponse{data=[]models.Commodity}
// @Router /public/commodities [get]
func (h *Handler) GetCommodities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	categoryID := c.Query("category_id")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	query := h.DB.Where("status = ? AND hide = ?", 1, 0).Preload("Category")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var total int64
	query.Model(&models.Commodity{}).Count(&total)

	var commodities []models.Commodity
	if err := query.Order("sort ASC, id DESC").Offset(offset).Limit(limit).Find(&commodities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch commodities"})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:  commodities,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// GetCommodity returns a single commodity by ID
// @Summary Get commodity
// @Description Get a single commodity by ID
// @Tags public
// @Produce json
// @Param id path int true "Commodity ID"
// @Success 200 {object} models.Commodity
// @Failure 404 {object} ErrorResponse
// @Router /public/commodity/{id} [get]
func (h *Handler) GetCommodity(c *gin.Context) {
	id := c.Param("id")

	var commodity models.Commodity
	if err := h.DB.Where("id = ? AND status = ? AND hide = ?", id, 1, 0).
		Preload("Category").First(&commodity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Commodity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch commodity"})
		return
	}

	c.JSON(http.StatusOK, commodity)
}

// GetPublicConfigs returns public configuration
// @Summary Get public configs
// @Description Get public system configuration
// @Tags public
// @Produce json
// @Success 200 {object} map[string]string
// @Router /public/configs [get]
func (h *Handler) GetPublicConfigs(c *gin.Context) {
	publicKeys := []string{
		"shop_name", "title", "description", "keywords",
		"notice", "service_qq", "service_url",
	}

	var configs []models.Config
	if err := h.DB.Where("key IN ?", publicKeys).Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch configs"})
		return
	}

	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	c.JSON(http.StatusOK, configMap)
}

// User handlers

// GetProfile returns the current user's profile
// @Summary Get user profile
// @Description Get the current user's profile information
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Router /user/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.User
	if err := h.DB.Preload("BusinessLevelInfo").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the current user's profile
// @Summary Update user profile
// @Description Update the current user's profile information
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateProfileRequest true "Profile update data"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Router /user/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update allowed fields
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.QQ != nil {
		user.QQ = req.QQ
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}
	if req.Nicename != nil {
		user.Nicename = req.Nicename
	}

	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserOrders returns the current user's orders
// @Summary Get user orders
// @Description Get the current user's order history with pagination
// @Tags user
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} PaginatedResponse{data=[]models.Order}
// @Router /user/orders [get]
func (h *Handler) GetUserOrders(c *gin.Context) {
	userID := c.GetUint("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var total int64
	h.DB.Model(&models.Order{}).Where("owner = ?", userID).Count(&total)

	var orders []models.Order
	if err := h.DB.Where("owner = ?", userID).
		Preload("Commodity").
		Preload("Payment").
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:  orders,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// CreateOrder creates a new order
// @Summary Create order
// @Description Create a new order for a commodity
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateOrderRequest true "Order data"
// @Success 201 {object} models.Order
// @Failure 400 {object} ErrorResponse
// @Router /user/orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate commodity exists and is available
	var commodity models.Commodity
	if err := h.DB.Where("id = ? AND status = ?", req.CommodityID, 1).First(&commodity).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Commodity not found or unavailable"})
		return
	}

	// Validate payment method exists
	var payment models.Payment
	if err := h.DB.Where("id = ? AND commodity = ?", req.PayID, 1).First(&payment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment method not found or unavailable"})
		return
	}

	// Generate trade number (simplified)
	tradeNo := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(int(userID))

	// Calculate amount based on user status
	amount := commodity.Price
	if userID > 0 {
		amount = commodity.UserPrice
	}

	if req.CardNum > 1 {
		amount *= float64(req.CardNum)
	}

	// Create order
	order := models.Order{
		Owner:        userID,
		UserID:       commodity.Owner,
		TradeNo:      tradeNo,
		Amount:       amount,
		CommodityID:  req.CommodityID,
		CardNum:      req.CardNum,
		PayID:        req.PayID,
		CreateIP:     c.ClientIP(),
		CreateDevice: req.CreateDevice,
		Contact:      req.Contact,
		Widget:       req.Widget,
		Status:       0, // Unpaid
	}

	if err := h.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Load associations for response
	h.DB.Preload("Commodity").Preload("Payment").First(&order, order.ID)

	c.JSON(http.StatusCreated, order)
}

// GetBalance returns the current user's balance
// @Summary Get user balance
// @Description Get the current user's balance and coin information
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} BalanceResponse
// @Router /user/balance [get]
func (h *Handler) GetBalance(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.User
	if err := h.DB.Select("balance, coin, integral, recharge, total_coin").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, BalanceResponse{
		Balance:   user.Balance,
		Coin:      user.Coin,
		Integral:  user.Integral,
		Recharge:  user.Recharge,
		TotalCoin: user.TotalCoin,
	})
}

// Admin handlers (simplified for now)

// GetDashboard returns admin dashboard data
func (h *Handler) GetDashboard(c *gin.Context) {
	var userCount, orderCount, commodityCount int64

	h.DB.Model(&models.User{}).Count(&userCount)
	h.DB.Model(&models.Order{}).Count(&orderCount)
	h.DB.Model(&models.Commodity{}).Count(&commodityCount)

	c.JSON(http.StatusOK, gin.H{
		"users":       userCount,
		"orders":      orderCount,
		"commodities": commodityCount,
	})
}

// GetUsers returns all users with pagination
func (h *Handler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var total int64
	h.DB.Model(&models.User{}).Count(&total)

	var users []models.User
	if err := h.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:  users,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// GetAllOrders returns all orders with pagination
func (h *Handler) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var total int64
	h.DB.Model(&models.Order{}).Count(&total)

	var orders []models.Order
	if err := h.DB.Preload("OwnerUser").Preload("Commodity").Preload("Payment").
		Offset(offset).Limit(limit).Order("id DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:  orders,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// Placeholder admin handlers for categories and commodities
func (h *Handler) GetAllCommodities(c *gin.Context)  { c.JSON(200, gin.H{"message": "TODO"}) }
func (h *Handler) CreateCommodity(c *gin.Context)    { c.JSON(200, gin.H{"message": "TODO"}) }
func (h *Handler) UpdateCommodity(c *gin.Context)    { c.JSON(200, gin.H{"message": "TODO"}) }
func (h *Handler) DeleteCommodity(c *gin.Context)    { c.JSON(200, gin.H{"message": "TODO"}) }
func (h *Handler) GetAllCategories(c *gin.Context)   { c.JSON(200, gin.H{"message": "TODO"}) }
func (h *Handler) CreateCategory(c *gin.Context)     { c.JSON(200, gin.H{"message": "TODO"}) }
func (h *Handler) UpdateCategory(c *gin.Context)     { c.JSON(200, gin.H{"message": "TODO"}) }
func (h *Handler) DeleteCategory(c *gin.Context)     { c.JSON(200, gin.H{"message": "TODO"}) }