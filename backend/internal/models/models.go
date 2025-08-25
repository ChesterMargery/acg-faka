package models

import (
	"time"
)

// User represents a user/customer in the system
type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Username         string    `json:"username" gorm:"unique;not null"`
	Email            string    `json:"email" gorm:"unique;not null"`
	Phone            string    `json:"phone"`
	QQ               string    `json:"qq"`
	Password         string    `json:"-" gorm:"not null"`
	Salt             string    `json:"-"`
	AppKey           string    `json:"app_key"`
	Avatar           string    `json:"avatar"`
	Balance          float64   `json:"balance" gorm:"default:0"`
	Coin             float64   `json:"coin" gorm:"default:0"`
	TotalCoin        float64   `json:"total_coin" gorm:"default:0"`
	Integral         int       `json:"integral" gorm:"default:0"`
	CreateTime       time.Time `json:"create_time"`
	LoginTime        time.Time `json:"login_time"`
	LastLoginTime    time.Time `json:"last_login_time"`
	LoginIP          string    `json:"login_ip"`
	LastLoginIP      string    `json:"last_login_ip"`
	PID              uint      `json:"pid" gorm:"default:0"`
	Status           int       `json:"status" gorm:"default:1"`
	BusinessLevel    int       `json:"business_level" gorm:"default:0"`
	Recharge         float64   `json:"recharge" gorm:"default:0"`
	Settlement       int       `json:"settlement" gorm:"default:0"`
	Nicename         string    `json:"nicename"`
	Alipay           string    `json:"alipay"`
	Wechat           string    `json:"wechat"`
}

// Card represents a virtual card/item for sale
type Card struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CommodityID  uint      `json:"commodity_id" gorm:"not null"`
	CreateTime   time.Time `json:"create_time"`
	Draft        string    `json:"draft"`
	OrderID      uint      `json:"order_id"`
	Owner        uint      `json:"owner"`
	PurchaseTime time.Time `json:"purchase_time"`
	Secret       string    `json:"secret"`
	Note         string    `json:"note"`
	Status       int       `json:"status" gorm:"default:0"` // 0: available, 1: sold
	Race         string    `json:"race"`
	
	// Relations
	OwnerUser User      `json:"owner_user,omitempty" gorm:"foreignKey:Owner"`
	Commodity Commodity `json:"commodity,omitempty" gorm:"foreignKey:CommodityID"`
	Order     Order     `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// Commodity represents a product/commodity
type Commodity struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null"`
	Description     string    `json:"description"`
	Cover           string    `json:"cover"`
	Price           float64   `json:"price" gorm:"not null"`
	VipPrice        float64   `json:"vip_price"`
	CategoryID      uint      `json:"category_id"`
	Owner           uint      `json:"owner"`
	Status          int       `json:"status" gorm:"default:1"`
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
	Sort            int       `json:"sort" gorm:"default:0"`
	IsAutoDelivery  bool      `json:"is_auto_delivery" gorm:"default:true"`
	
	// Relations
	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	OwnerUser User    `json:"owner_user,omitempty" gorm:"foreignKey:Owner"`
	Cards    []Card  `json:"cards,omitempty" gorm:"foreignKey:CommodityID"`
}

// Category represents a product category
type Category struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	Icon       string    `json:"icon"`
	Sort       int       `json:"sort" gorm:"default:0"`
	Owner      uint      `json:"owner"`
	Status     int       `json:"status" gorm:"default:1"`
	CreateTime time.Time `json:"create_time"`
	
	// Relations
	OwnerUser   User        `json:"owner_user,omitempty" gorm:"foreignKey:Owner"`
	Commodities []Commodity `json:"commodities,omitempty" gorm:"foreignKey:CategoryID"`
}

// Order represents a purchase order
type Order struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TradeNo      string    `json:"trade_no" gorm:"unique;not null"`
	UserID       uint      `json:"user_id"`
	CommodityID  uint      `json:"commodity_id"`
	Quantity     int       `json:"quantity" gorm:"default:1"`
	Amount       float64   `json:"amount" gorm:"not null"`
	Status       int       `json:"status" gorm:"default:0"` // 0: pending, 1: paid, 2: cancelled
	PayMethod    string    `json:"pay_method"`
	CreateTime   time.Time `json:"create_time"`
	PayTime      time.Time `json:"pay_time"`
	Email        string    `json:"email"`
	Race         string    `json:"race"`
	
	// Relations
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Commodity Commodity `json:"commodity,omitempty" gorm:"foreignKey:CommodityID"`
	Cards     []Card    `json:"cards,omitempty" gorm:"foreignKey:OrderID"`
}

// Manage represents an admin user
type Manage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Username   string    `json:"username" gorm:"unique;not null"`
	Email      string    `json:"email" gorm:"unique;not null"`
	Password   string    `json:"-" gorm:"not null"`
	Salt       string    `json:"-"`
	Status     int       `json:"status" gorm:"default:1"`
	CreateTime time.Time `json:"create_time"`
	LoginTime  time.Time `json:"login_time"`
	LoginIP    string    `json:"login_ip"`
}

// Config represents system configuration
type Config struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Key   string `json:"key" gorm:"unique;not null"`
	Value string `json:"value"`
	Type  string `json:"type" gorm:"default:string"`
}

// Bill represents transaction records
type Bill struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Amount     float64   `json:"amount" gorm:"not null"`
	Type       int       `json:"type" gorm:"not null"` // 1: income, 2: expense
	Note       string    `json:"note"`
	CreateTime time.Time `json:"create_time"`
	
	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Cash represents withdrawal records
type Cash struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Amount     float64   `json:"amount" gorm:"not null"`
	Cost       float64   `json:"cost" gorm:"default:0"`
	Type       int       `json:"type" gorm:"default:1"`
	Card       int       `json:"card" gorm:"default:1"`
	Status     int       `json:"status" gorm:"default:0"` // 0: pending, 1: approved, 2: rejected
	CreateTime time.Time `json:"create_time"`
	
	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}