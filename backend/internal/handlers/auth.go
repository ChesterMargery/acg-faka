package handlers

import (
	"acg-faka/backend/internal/database"
	"acg-faka/backend/internal/middleware"
	"acg-faka/backend/internal/models"
	"acg-faka/backend/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct{}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	QQ       string `json:"qq"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// UserLogin handles user login
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account disabled"})
		return
	}

	// Update login time and IP
	user.LastLoginTime = user.LoginTime
	user.LoginTime = time.Now()
	user.LastLoginIP = user.LoginIP
	user.LoginIP = c.ClientIP()
	database.DB.Save(&user)

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	// Remove sensitive data
	user.Password = ""
	user.Salt = ""

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

// AdminLogin handles admin login
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Manage
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, admin.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if admin.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account disabled"})
		return
	}

	// Update login time and IP
	admin.LoginTime = time.Now()
	admin.LoginIP = c.ClientIP()
	database.DB.Save(&admin)

	// Generate token
	token, err := middleware.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	// Remove sensitive data
	admin.Password = ""
	admin.Salt = ""

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  admin,
	})
}

// UserRegister handles user registration
func (h *AuthHandler) UserRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username or email already exists
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	// Create user
	user := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Phone:      req.Phone,
		QQ:         req.QQ,
		Password:   hashedPassword,
		Salt:       utils.GenerateSalt(),
		AppKey:     utils.GenerateAppKey(),
		CreateTime: time.Now(),
		LoginTime:  time.Now(),
		LoginIP:    c.ClientIP(),
		Status:     1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	// Remove sensitive data
	user.Password = ""
	user.Salt = ""

	c.JSON(http.StatusCreated, LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetProfile returns user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")

	if userType == "admin" {
		var admin models.Manage
		if err := database.DB.First(&admin, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
			return
		}
		admin.Password = ""
		admin.Salt = ""
		c.JSON(http.StatusOK, admin)
	} else {
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		user.Password = ""
		user.Salt = ""
		c.JSON(http.StatusOK, user)
	}
}