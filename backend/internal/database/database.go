package database

import (
	"acg-faka/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

// InitDatabase initializes the SQLite database connection
func InitDatabase(dbPath string) error {
	var err error
	
	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	
	// Open database connection
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}
	
	// Auto migrate the schema
	return autoMigrate()
}

// autoMigrate migrates all models
func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Card{},
		&models.Commodity{},
		&models.Category{},
		&models.Order{},
		&models.Manage{},
		&models.Config{},
		&models.Bill{},
		&models.Cash{},
	)
}

// SeedDatabase seeds initial data
func SeedDatabase() error {
	// Create default admin user
	var count int64
	DB.Model(&models.Manage{}).Count(&count)
	if count == 0 {
		// Hash password using the same method as in handlers
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		
		admin := models.Manage{
			Username:   "admin",
			Email:      "admin@admin.com",
			Password:   string(hashedPassword),
			Salt:       "acg-faka",
			Status:     1,
			CreateTime: time.Now(),
		}
		if err := DB.Create(&admin).Error; err != nil {
			return err
		}
	}
	
	// Create default config values
	defaultConfigs := []models.Config{
		{Key: "site_name", Value: "ACG-FAKA", Type: "string"},
		{Key: "site_description", Value: "Virtual Card Store", Type: "string"},
		{Key: "site_logo", Value: "/favicon.ico", Type: "string"},
		{Key: "cash_min", Value: "10", Type: "float"},
		{Key: "cash_cost", Value: "1", Type: "float"},
		{Key: "trade_verification", Value: "0", Type: "int"},
		{Key: "cash_type_balance", Value: "1", Type: "int"},
	}
	
	for _, config := range defaultConfigs {
		var existingConfig models.Config
		if err := DB.Where("key = ?", config.Key).First(&existingConfig).Error; err == gorm.ErrRecordNotFound {
			if err := DB.Create(&config).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}