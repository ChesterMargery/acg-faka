package handlers

import "acg-faka/internal/models"

// Request/Response DTOs

type RegisterRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=32"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone"`
	QQ       *string `json:"qq"`
	Password string  `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type AdminAuthResponse struct {
	Token   string          `json:"token"`
	Manager *models.Manager `json:"manager"`
}

type UpdateProfileRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone"`
	QQ       *string `json:"qq"`
	Avatar   *string `json:"avatar"`
	Nicename *string `json:"nicename"`
}

type CreateOrderRequest struct {
	CommodityID  uint    `json:"commodity_id" binding:"required"`
	PayID        uint    `json:"pay_id" binding:"required"`
	CardNum      int     `json:"card_num" binding:"min=1"`
	Contact      *string `json:"contact"`
	CreateDevice int     `json:"create_device"` // 0=pc, 1=android, 2=ios, 3=ipad
	Widget       *string `json:"widget"`
}

type BalanceResponse struct {
	Balance   float64 `json:"balance"`
	Coin      float64 `json:"coin"`
	Integral  int     `json:"integral"`
	Recharge  float64 `json:"recharge"`
	TotalCoin float64 `json:"total_coin"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}