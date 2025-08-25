package database

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(databaseURL string) *gorm.DB {
	// Ensure the data directory exists
	dataDir := filepath.Dir(databaseURL)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Open database connection
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")
	return db
}

func RunMigrations(databaseURL string) error {
	// For now, we'll use a simple approach to run the SQL migration
	// In a production system, you might want to use golang-migrate
	log.Println("Running database migrations...")
	
	// Check if the database file exists and has tables
	// If not, run the initial migration
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Check if users table exists
	var count int
	err = sqlDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if err != nil || count == 0 {
		// Run initial migration
		migrationPath := "migrations/001_initial_schema.sql"
		if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
			// Try relative path from backend directory
			migrationPath = "../migrations/001_initial_schema.sql"
			if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
				log.Printf("Warning: Could not find migration file: %v", err)
				return nil // Don't fail if migration file is not found
			}
		}

		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			log.Printf("Warning: Could not read migration file: %v", err)
			return nil // Don't fail if migration file is not found
		}

		if err := sqlDB.Ping(); err != nil {
			return err
		}

		// Execute migration
		if _, err := sqlDB.Exec(string(migrationSQL)); err != nil {
			log.Printf("Warning: Migration failed: %v", err)
			return nil // Don't fail on migration errors for now
		}

		log.Println("Database migrations completed successfully")
	} else {
		log.Println("Database already initialized")
	}

	return nil
}