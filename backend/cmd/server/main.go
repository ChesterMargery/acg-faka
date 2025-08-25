package main

import (
	"acg-faka/internal/config"
	"acg-faka/internal/database"
	"acg-faka/internal/handlers"
	"acg-faka/internal/middleware"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title ACG-FAKA API
// @version 2.0
// @description A modern virtual card selling system API built with Go
// @termsOfService http://swagger.io/terms/

// @contact.name ACG-FAKA Support
// @contact.url https://github.com/ChesterMargery/acg-faka
// @contact.email support@acg-faka.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db := database.Init(cfg.DatabaseURL)

	// Run migrations
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// Setup CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"} // Vue dev servers
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(config))

	// Setup middleware
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// Setup routes
	api := router.Group("/api/v1")
	handlers.SetupRoutes(api, db)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "ACG-FAKA API is running"})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}