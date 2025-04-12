package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/logger"
	"github.com/jaimesHub/golang-todo-app/internal/services/auth"
)

// LoggingMiddleware creates a middleware for logging requests and responses
func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Append query string if exists
		if raw != "" {
			path = path + "?" + raw
		}

		// Log request details
		log.Info("Request processed",
			map[string]interface{}{
				"status":     statusCode,
				"method":     method,
				"path":       path,
				"latency":    latency,
				"client_ip":  clientIP,
				"user_agent": userAgent,
			})
	}
}

// AuthMiddleware verifies JWT token
func AuthMiddleware(jwtService *auth.JWTService, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("Missing Authorization header", map[string]interface{}{
				"path":      c.Request.URL.Path,
				"client_ip": c.ClientIP(),
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn("Invalid Authorization header format", map[string]interface{}{
				"path":      c.Request.URL.Path,
				"client_ip": c.ClientIP(),
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Extract token
		tokenString := parts[1]

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Warn("Invalid or expired token", map[string]interface{}{
				"path":      c.Request.URL.Path,
				"client_ip": c.ClientIP(),
				"error":     err.Error(),
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", claims.UserID)

		log.Debug("User authenticated", map[string]interface{}{
			"user_id":   claims.UserID,
			"path":      c.Request.URL.Path,
			"client_ip": c.ClientIP(),
		})

		c.Next()
	}
}

// Setup configures middleware for the router
func Setup(router *gin.Engine, cfg *config.Config, jwtService *auth.JWTService, log *logger.Logger) {
	// Add CORS middleware
	router.Use(corsMiddleware())

	// Add logging middleware
	router.Use(LoggingMiddleware(log))

	// Add recovery middleware with logging
	router.Use(RecoveryMiddleware(log))
}

// corsMiddleware handles Cross-Origin Resource Sharing
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered", map[string]interface{}{
					"error":     err,
					"path":      c.Request.URL.Path,
					"method":    c.Request.Method,
					"client_ip": c.ClientIP(),
				})
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
