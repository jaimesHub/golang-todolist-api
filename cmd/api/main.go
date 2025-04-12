package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/database"
	"github.com/jaimesHub/golang-todo-app/internal/middleware"
	"github.com/jaimesHub/golang-todo-app/internal/routes"
	"github.com/jaimesHub/golang-todo-app/internal/services/redis"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis client
	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		defer redisClient.Close()
	}

	// Initialize Gin router
	router := gin.Default()

	// Apply middleware
	middleware.Setup(router, cfg)

	// Register routes
	routes.Register(router, db, redisClient)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
