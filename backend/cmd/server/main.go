package main

import (
	"acg-faka/backend/internal/database"
	"acg-faka/backend/internal/handlers"
	"acg-faka/backend/internal/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./acg-faka.db"
	}

	if err := database.InitDatabase(dbPath); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Seed database with default data
	if err := database.SeedDatabase(); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORSMiddleware())

	// Initialize handlers
	authHandler := &handlers.AuthHandler{}
	categoryHandler := &handlers.CategoryHandler{}
	commodityHandler := &handlers.CommodityHandler{}
	cardHandler := &handlers.CardHandler{}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "ACG-FAKA API is running"})
	})

	// Public routes
	api := r.Group("/api")
	{
		// Authentication
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.UserLogin)
			auth.POST("/register", authHandler.UserRegister)
			auth.POST("/admin/login", authHandler.AdminLogin)
		}

		// Public commodity and category endpoints
		public := api.Group("/public")
		{
			public.GET("/categories", categoryHandler.GetCategories)
			public.GET("/categories/:id", categoryHandler.GetCategory)
			public.GET("/commodities", commodityHandler.GetCommodities)
			public.GET("/commodities/:id", commodityHandler.GetCommodity)
		}
	}

	// Protected user routes
	user := api.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/profile", authHandler.GetProfile)
		
		// User categories
		user.GET("/categories", categoryHandler.GetCategories)
		user.POST("/categories", categoryHandler.CreateCategory)
		user.PUT("/categories/:id", categoryHandler.UpdateCategory)
		user.DELETE("/categories/:id", categoryHandler.DeleteCategory)

		// User commodities
		user.GET("/commodities", commodityHandler.GetCommodities)
		user.POST("/commodities", commodityHandler.CreateCommodity)
		user.PUT("/commodities/:id", commodityHandler.UpdateCommodity)
		user.DELETE("/commodities/:id", commodityHandler.DeleteCommodity)

		// User cards
		user.GET("/cards", cardHandler.GetCards)
		user.GET("/cards/:id", cardHandler.GetCard)
		user.POST("/cards", cardHandler.CreateCard)
		user.POST("/cards/batch", cardHandler.BatchCreateCards)
		user.PUT("/cards/:id", cardHandler.UpdateCard)
		user.DELETE("/cards/:id", cardHandler.DeleteCard)
		user.DELETE("/cards/batch", cardHandler.BatchDeleteCards)
	}

	// Protected admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		admin.GET("/profile", authHandler.GetProfile)
		
		// Admin categories
		admin.GET("/categories", categoryHandler.GetCategories)
		admin.POST("/categories", categoryHandler.CreateCategory)
		admin.PUT("/categories/:id", categoryHandler.UpdateCategory)
		admin.DELETE("/categories/:id", categoryHandler.DeleteCategory)

		// Admin commodities
		admin.GET("/commodities", commodityHandler.GetCommodities)
		admin.POST("/commodities", commodityHandler.CreateCommodity)
		admin.PUT("/commodities/:id", commodityHandler.UpdateCommodity)
		admin.DELETE("/commodities/:id", commodityHandler.DeleteCommodity)

		// Admin cards
		admin.GET("/cards", cardHandler.GetCards)
		admin.GET("/cards/:id", cardHandler.GetCard)
		admin.POST("/cards", cardHandler.CreateCard)
		admin.POST("/cards/batch", cardHandler.BatchCreateCards)
		admin.PUT("/cards/:id", cardHandler.UpdateCard)
		admin.DELETE("/cards/:id", cardHandler.DeleteCard)
		admin.DELETE("/cards/batch", cardHandler.BatchDeleteCards)
	}

	// Serve static files for frontend
	r.Static("/static", "./frontend/dist")
	r.StaticFile("/", "./frontend/dist/index.html")
	r.StaticFile("/admin", "./frontend/dist/admin.html")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}