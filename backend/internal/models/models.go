package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Username        string     `json:"username" gorm:"uniqueIndex;not null"`
	Email           *string    `json:"email" gorm:"uniqueIndex"`
	Phone           *string    `json:"phone" gorm:"uniqueIndex"`
	QQ              *string    `json:"qq"`
	Password        string     `json:"-" gorm:"not null"`
	Salt            string     `json:"-" gorm:"not null"`
	AppKey          string     `json:"app_key" gorm:"not null"`
	Avatar          *string    `json:"avatar"`
	Balance         float64    `json:"balance" gorm:"type:decimal(14,2);default:0.00"`
	Coin            float64    `json:"coin" gorm:"type:decimal(14,2);default:0.00"`
	Integral        int        `json:"integral" gorm:"default:0"`
	CreateTime      time.Time  `json:"create_time" gorm:"autoCreateTime"`
	LoginTime       *time.Time `json:"login_time"`
	LastLoginTime   *time.Time `json:"last_login_time"`
	LoginIP         *string    `json:"login_ip"`
	LastLoginIP     *string    `json:"last_login_ip"`
	PID             uint       `json:"pid" gorm:"column:pid;default:0"`
	Recharge        float64    `json:"recharge" gorm:"type:decimal(14,2);default:0.00"`
	TotalCoin       float64    `json:"total_coin" gorm:"type:decimal(14,2);default:0.00"`
	Status          int        `json:"status" gorm:"default:0"` // 0=banned, 1=normal
	BusinessLevel   *uint      `json:"business_level"`
	Nicename        *string    `json:"nicename"`
	Alipay          *string    `json:"alipay"`
	Wechat          *string    `json:"wechat"`
	Settlement      int        `json:"settlement" gorm:"default:0"` // 0=alipay, 1=wechat

	// Associations
	Parent         *User          `json:"parent,omitempty" gorm:"foreignKey:PID"`
	BusinessLevelInfo *BusinessLevel `json:"business_level_info,omitempty" gorm:"foreignKey:BusinessLevel"`
}

// Business represents a business/shop
type Business struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	ShopName      *string        `json:"shop_name"`
	Title         *string        `json:"title"`
	Notice        *string        `json:"notice"`
	ServiceQQ     *string        `json:"service_qq"`
	ServiceURL    *string        `json:"service_url"`
	Subdomain     *string        `json:"subdomain" gorm:"uniqueIndex"`
	Topdomain     *string        `json:"topdomain" gorm:"uniqueIndex"`
	MasterDisplay int            `json:"master_display" gorm:"default:0"`
	CreateTime    time.Time      `json:"create_time" gorm:"autoCreateTime"`

	// Associations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// BusinessLevel represents business level configuration
type BusinessLevel struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	Name       string  `json:"name" gorm:"not null"`
	Icon       *string `json:"icon"`
	Cost       float64 `json:"cost" gorm:"type:decimal(4,2);default:0.00"`
	Accrual    float64 `json:"accrual" gorm:"type:decimal(4,2);default:0.00"`
	Substation int     `json:"substation" gorm:"default:0"`
	TopDomain  int     `json:"top_domain" gorm:"default:0"`
	Price      float64 `json:"price" gorm:"type:decimal(10,2);default:0.00"`
	Supplier   int     `json:"supplier" gorm:"default:1"`
}

// Category represents a product category
type Category struct {
	ID              uint    `json:"id" gorm:"primaryKey"`
	Name            string  `json:"name" gorm:"not null"`
	Sort            int     `json:"sort" gorm:"default:0"`
	CreateTime      time.Time `json:"create_time" gorm:"autoCreateTime"`
	Owner           uint    `json:"owner" gorm:"default:0"`
	Icon            *string `json:"icon"`
	Status          int     `json:"status" gorm:"default:0"`
	Hide            int     `json:"hide" gorm:"default:0"`
	UserLevelConfig *string `json:"user_level_config"`
}

// Commodity represents a product/commodity
type Commodity struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	CategoryID         uint           `json:"category_id" gorm:"not null"`
	Name               string         `json:"name" gorm:"not null"`
	Description        *string        `json:"description"`
	Cover              *string        `json:"cover"`
	FactoryPrice       float64        `json:"factory_price" gorm:"type:decimal(10,2);default:0.00"`
	Price              float64        `json:"price" gorm:"type:decimal(10,2);default:0.00"`
	UserPrice          float64        `json:"user_price" gorm:"type:decimal(10,2);default:0.00"`
	Status             int            `json:"status" gorm:"default:0"`
	Owner              uint           `json:"owner" gorm:"default:0"`
	CreateTime         time.Time      `json:"create_time" gorm:"autoCreateTime"`
	APIStatus          int            `json:"api_status" gorm:"default:0"`
	Code               string         `json:"code" gorm:"uniqueIndex;not null"`
	DeliveryWay        int            `json:"delivery_way" gorm:"default:0"`
	DeliveryAutoMode   int            `json:"delivery_auto_mode" gorm:"default:0"`
	DeliveryMessage    *string        `json:"delivery_message"`
	ContactType        int            `json:"contact_type" gorm:"default:0"`
	PasswordStatus     int            `json:"password_status" gorm:"default:0"`
	Sort               int            `json:"sort" gorm:"default:0"`
	Coupon             int            `json:"coupon" gorm:"default:0"`
	SharedID           *uint          `json:"shared_id"`
	SharedCode         *string        `json:"shared_code"`
	SharedPremium      *float64       `json:"shared_premium" gorm:"type:decimal(10,2);default:0.00"`
	SharedPremiumType  *int           `json:"shared_premium_type" gorm:"default:0"`
	SeckillStatus      int            `json:"seckill_status" gorm:"default:0"`
	SeckillStartTime   *time.Time     `json:"seckill_start_time"`
	SeckillEndTime     *time.Time     `json:"seckill_end_time"`
	DraftStatus        int            `json:"draft_status" gorm:"default:0"`
	DraftPremium       *float64       `json:"draft_premium" gorm:"type:decimal(10,2);default:0.00"`
	InventoryHidden    int            `json:"inventory_hidden" gorm:"default:0"`
	LeaveMessage       *string        `json:"leave_message"`
	Recommend          *int           `json:"recommend" gorm:"default:0"`
	SendEmail          int            `json:"send_email" gorm:"default:0"`
	OnlyUser           int            `json:"only_user" gorm:"default:0"`
	PurchaseCount      int            `json:"purchase_count" gorm:"default:0"`
	Widget             *string        `json:"widget"`
	LevelPrice         *string        `json:"level_price"`
	LevelDisable       int            `json:"level_disable" gorm:"default:0"`
	Minimum            int            `json:"minimum" gorm:"default:0"`
	Maximum            int            `json:"maximum" gorm:"default:0"`
	SharedSync         int            `json:"shared_sync" gorm:"default:0"`
	Config             *string        `json:"config"`
	Hide               int            `json:"hide" gorm:"default:0"`
	InventorySync      int            `json:"inventory_sync" gorm:"default:0"`

	// Associations
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Cards    []Card    `json:"cards,omitempty" gorm:"foreignKey:CommodityID"`
}

// Card represents a card/key
type Card struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Owner        uint           `json:"owner" gorm:"default:0"`
	CommodityID  uint           `json:"commodity_id" gorm:"not null"`
	Draft        *string        `json:"draft"`
	Secret       string         `json:"secret" gorm:"not null"`
	CreateTime   time.Time      `json:"create_time" gorm:"autoCreateTime"`
	PurchaseTime *time.Time     `json:"purchase_time"`
	OrderID      *uint          `json:"order_id"`
	Status       int            `json:"status" gorm:"default:0"`
	Note         *string        `json:"note"`
	Race         *string        `json:"race"`

	// Associations
	Commodity *Commodity `json:"commodity,omitempty" gorm:"foreignKey:CommodityID"`
	Order     *Order     `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// Order represents an order
type Order struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Owner          uint           `json:"owner" gorm:"default:0"`
	UserID         uint           `json:"user_id" gorm:"default:0"`
	TradeNo        string         `json:"trade_no" gorm:"uniqueIndex;not null"`
	Amount         float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	CommodityID    uint           `json:"commodity_id" gorm:"not null"`
	CardID         *uint          `json:"card_id"`
	CardNum        int            `json:"card_num" gorm:"default:0"`
	PayID          uint           `json:"pay_id" gorm:"not null"`
	CreateTime     time.Time      `json:"create_time" gorm:"autoCreateTime"`
	CreateIP       string         `json:"create_ip" gorm:"not null"`
	CreateDevice   int            `json:"create_device" gorm:"not null"`
	PayTime        *time.Time     `json:"pay_time"`
	Status         int            `json:"status" gorm:"default:0"`
	Secret         *string        `json:"secret"`
	Password       *string        `json:"password"`
	Contact        *string        `json:"contact"`
	DeliveryStatus int            `json:"delivery_status" gorm:"default:0"`
	PayURL         *string        `json:"pay_url"`
	CouponID       *uint          `json:"coupon_id"`
	Cost           float64        `json:"cost" gorm:"type:decimal(10,2);default:0.00"`
	FromUserID     *uint          `json:"from_user_id"`
	Premium        *float64       `json:"premium" gorm:"type:decimal(10,2);default:0.00"`
	Widget         *string        `json:"widget"`
	Rent           float64        `json:"rent" gorm:"type:decimal(10,2);default:0.00"`
	Race           *string        `json:"race"`
	Rebate         *float64       `json:"rebate" gorm:"type:decimal(10,2);default:0.00"`
	PayCost        *float64       `json:"pay_cost" gorm:"type:decimal(10,2);default:0.00"`
	RequestNo      *string        `json:"request_no" gorm:"uniqueIndex"`

	// Associations
	OwnerUser *User      `json:"owner_user,omitempty" gorm:"foreignKey:Owner"`
	User      *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Commodity *Commodity `json:"commodity,omitempty" gorm:"foreignKey:CommodityID"`
	Card      *Card      `json:"card,omitempty" gorm:"foreignKey:CardID"`
	Payment   *Payment   `json:"payment,omitempty" gorm:"foreignKey:PayID"`
	FromUser  *User      `json:"from_user,omitempty" gorm:"foreignKey:FromUserID"`
	Coupon    *Coupon    `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
}

// Payment represents a payment method
type Payment struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"not null"`
	Icon       *string        `json:"icon"`
	Code       string         `json:"code" gorm:"not null"`
	Commodity  int            `json:"commodity" gorm:"default:0"`
	Recharge   int            `json:"recharge" gorm:"default:0"`
	CreateTime time.Time      `json:"create_time" gorm:"autoCreateTime"`
	Handle     string         `json:"handle" gorm:"not null"`
	Sort       int            `json:"sort" gorm:"default:0"`
	Equipment  int            `json:"equipment" gorm:"default:0"`
	Cost       *float64       `json:"cost" gorm:"type:decimal(10,3);default:0.000"`
	CostType   *int           `json:"cost_type" gorm:"default:0"`
}

// Config represents system configuration
type Config struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Key   string `json:"key" gorm:"uniqueIndex;not null"`
	Value string `json:"value" gorm:"not null"`
}

// Coupon represents a coupon
type Coupon struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null"`
	CommodityID uint           `json:"commodity_id" gorm:"not null"`
	Owner       uint           `json:"owner" gorm:"default:0"`
	CreateTime  time.Time      `json:"create_time" gorm:"autoCreateTime"`
	ExpireTime  *time.Time     `json:"expire_time"`
	ServiceTime *time.Time     `json:"service_time"`
	Money       float64        `json:"money" gorm:"type:decimal(10,2);not null"`
	Status      int            `json:"status" gorm:"default:0"`
	TradeNo     *string        `json:"trade_no"`
	Note        *string        `json:"note"`
	Mode        *int           `json:"mode" gorm:"default:0"`
	CategoryID  *uint          `json:"category_id" gorm:"default:0"`
	Life        int            `json:"life" gorm:"default:1"`
	UseLife     int            `json:"use_life" gorm:"default:0"`
	Race        *string        `json:"race"`

	// Associations
	Commodity *Commodity `json:"commodity,omitempty" gorm:"foreignKey:CommodityID"`
}

// Manager represents an admin manager
type Manager struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Email            string         `json:"email" gorm:"uniqueIndex;not null"`
	Password         string         `json:"-" gorm:"not null"`
	SecurityPassword *string        `json:"-"`
	Nickname         *string        `json:"nickname"`
	Salt             string         `json:"-" gorm:"not null"`
	Avatar           *string        `json:"avatar"`
	Status           int            `json:"status" gorm:"default:0"`
	Type             int            `json:"type" gorm:"default:0"`
	CreateTime       time.Time      `json:"create_time" gorm:"autoCreateTime"`
	LoginTime        *time.Time     `json:"login_time"`
	LastLoginTime    *time.Time     `json:"last_login_time"`
	LoginIP          *string        `json:"login_ip"`
	LastLoginIP      *string        `json:"last_login_ip"`
	Note             *string        `json:"note"`
}

// UserGroup represents user membership levels
type UserGroup struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name" gorm:"not null"`
	Icon     *string `json:"icon"`
	Discount float64 `json:"discount" gorm:"type:decimal(4,2);not null"`
	Cost     float64 `json:"cost" gorm:"type:decimal(4,2);default:0.00"`
	Recharge float64 `json:"recharge" gorm:"type:decimal(14,2);uniqueIndex;not null"`
}