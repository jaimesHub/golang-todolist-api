package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/database"
	"github.com/jaimesHub/golang-todo-app/internal/middleware"
	"github.com/jaimesHub/golang-todo-app/internal/routes"
	"github.com/jaimesHub/golang-todo-app/internal/services/auth"
	"github.com/jaimesHub/golang-todo-app/internal/services/queue"
	"github.com/jaimesHub/golang-todo-app/internal/services/redis"
	"github.com/jaimesHub/golang-todo-app/internal/services/worker"
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

	// Initialize queue
	taskQueue, err := queue.NewQueue(cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to initialize task queue: %v", err)
	} else {
		defer taskQueue.Close()
	}

	// Initialize worker
	taskWorker, err := worker.NewWorker(cfg.Redis, "tasks")
	if err != nil {
		log.Printf("Warning: Failed to initialize task worker: %v", err)
	} else {
		defer taskWorker.Close()

		// Register task handlers
		taskWorker.RegisterHandler("email_notification", func(task *queue.Task) error {
			log.Printf("Processing email notification task: %v", task.Data)
			// Simulate work
			time.Sleep(1 * time.Second)
			return nil
		})

		taskWorker.RegisterHandler("task_reminder", func(task *queue.Task) error {
			log.Printf("Processing task reminder: %v", task.Data)
			// Simulate work
			time.Sleep(1 * time.Second)
			return nil
		})

		// Start worker
		taskWorker.Start()
	}

	// Initialize JWT service
	jwtService := auth.NewJWTService(&cfg.JWT)

	// Initialize Gin router
	router := gin.Default()

	// Apply middleware
	middleware.Setup(router, cfg, jwtService)

	// Register routes
	routes.Register(router, db, redisClient, cfg)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
