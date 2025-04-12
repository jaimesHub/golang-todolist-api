package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/services/auth"
)

// AuthMiddleware verifies JWT token
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Extract token
		tokenString := parts[1]

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

// Setup configures middleware for the router
func Setup(router *gin.Engine, cfg *config.Config, jwtService *auth.JWTService) {
	// Add CORS middleware
	router.Use(corsMiddleware())

	// Add logging middleware
	router.Use(loggingMiddleware())

	// Add recovery middleware
	router.Use(gin.Recovery())
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

// loggingMiddleware logs request information
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Log request details
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Log using structured logging
		// logger.Info("Request processed",
		//    "status", statusCode,
		//    "method", method,
		//    "path", path,
		//    "client_ip", clientIP)
	}
}
