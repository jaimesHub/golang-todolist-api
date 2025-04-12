package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/database"
	"github.com/jaimesHub/golang-todo-app/internal/logger"
	"github.com/jaimesHub/golang-todo-app/internal/middleware"
	"github.com/jaimesHub/golang-todo-app/internal/monitoring"
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

	// Initialize logger
	appLogger, err := logger.NewLogger(cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	appLogger.Info("Starting application", map[string]interface{}{
		"environment": gin.Mode(),
		"version":     "1.0.0",
	})

	// Initialize database connection
	db, err := database.Connect(cfg.Database)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", map[string]interface{}{"error": err.Error()})
	}
	appLogger.Info("Connected to database")

	// Run migrations
	if err := database.Migrate(db); err != nil {
		appLogger.Fatal("Failed to run migrations", map[string]interface{}{"error": err.Error()})
	}
	appLogger.Info("Database migrations completed")

	// Initialize Redis client
	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		appLogger.Warn("Failed to connect to Redis", map[string]interface{}{"error": err.Error()})
	} else {
		defer redisClient.Close()
		appLogger.Info("Connected to Redis")
	}

	// Initialize queue
	taskQueue, err := queue.NewQueue(cfg.Redis)
	if err != nil {
		appLogger.Warn("Failed to initialize task queue", map[string]interface{}{"error": err.Error()})
	} else {
		defer taskQueue.Close()
		appLogger.Info("Task queue initialized")
	}

	// Initialize worker
	taskWorker, err := worker.NewWorker(cfg.Redis, "tasks")
	if err != nil {
		appLogger.Warn("Failed to initialize task worker", map[string]interface{}{"error": err.Error()})
	} else {
		defer taskWorker.Close()

		// Register task handlers
		taskWorker.RegisterHandler("email_notification", func(task *queue.Task) error {
			appLogger.Info("Processing email notification task", map[string]interface{}{"task_id": task.ID, "data": task.Data})
			// Simulate work
			time.Sleep(1 * time.Second)
			return nil
		})

		taskWorker.RegisterHandler("task_reminder", func(task *queue.Task) error {
			appLogger.Info("Processing task reminder", map[string]interface{}{"task_id": task.ID, "data": task.Data})
			// Simulate work
			time.Sleep(1 * time.Second)
			return nil
		})

		// Start worker
		taskWorker.Start()
		appLogger.Info("Task worker started")
	}

	// Initialize JWT service
	jwtService := auth.NewJWTService(&cfg.JWT)

	// Initialize Gin router
	router := gin.Default()

	// Apply middleware
	middleware.Setup(router, cfg, jwtService, appLogger)
	router.Use(monitoring.MetricsMiddleware(appLogger))

	// Setup health check and metrics endpoints
	monitoring.SetupHealthCheck(router, appLogger)

	// Register routes
	routes.Register(router, db, redisClient, cfg)
	appLogger.Info("Routes registered")

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	appLogger.Info("Starting server", map[string]interface{}{"address": serverAddr})
	if err := router.Run(serverAddr); err != nil {
		appLogger.Fatal("Failed to start server", map[string]interface{}{"error": err.Error()})
	}
}
