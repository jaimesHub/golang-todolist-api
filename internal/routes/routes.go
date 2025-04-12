package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/handlers"
	"github.com/jaimesHub/golang-todo-app/internal/middleware"
	"github.com/jaimesHub/golang-todo-app/internal/services"
	"github.com/jaimesHub/golang-todo-app/internal/services/auth"
	redisService "github.com/jaimesHub/golang-todo-app/internal/services/redis"
	"gorm.io/gorm"
)

// Register sets up all API routes
func Register(router *gin.Engine, db *gorm.DB, redisClient *redisService.Client, cfg *config.Config) {
	// Create services
	userService := services.NewUserService(db)
	jwtService := auth.NewJWTService(&cfg.JWT)

	// Create handlers with dependencies
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService, jwtService)
	taskHandler := handlers.NewTaskHandler(db)

	// Setup middleware
	middleware.Setup(router, cfg, jwtService)

	// Health check route
	router.GET("/health", handlers.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes - no authentication required
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes - authentication required
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetProfile)
				users.PUT("/me", userHandler.UpdateProfile)
				users.GET("/activities", userHandler.GetActivities)
			}

			// Task routes
			tasks := protected.Group("/tasks")
			{
				tasks.POST("/", taskHandler.Create)
				tasks.GET("/", taskHandler.List)
				tasks.GET("/:id", taskHandler.GetByID)
				tasks.PUT("/:id", taskHandler.Update)
				tasks.DELETE("/:id", taskHandler.Delete)
			}
		}
	}
}
