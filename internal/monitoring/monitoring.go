package monitoring

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaimesHub/golang-todo-app/internal/logger"
)

// HealthCheck handles the health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// MetricsMiddleware collects basic metrics for monitoring
func MetricsMiddleware(log *logger.Logger) gin.HandlerFunc {
	// In a real application, you might want to use a metrics library like Prometheus
	// For simplicity, we'll just log some basic metrics

	// Initialize counters
	requestCount := 0
	responseStatusCount := make(map[int]int)

	return func(c *gin.Context) {
		// Increment request counter
		requestCount++

		// Process request
		c.Next()

		// Record response status
		status := c.Writer.Status()
		responseStatusCount[status]++

		// Log metrics periodically (every 100 requests)
		if requestCount%100 == 0 {
			log.Info("Metrics update", map[string]interface{}{
				"total_requests": requestCount,
				"status_counts":  responseStatusCount,
			})
		}
	}
}

// SetupHealthCheck registers health check endpoints
func SetupHealthCheck(router *gin.Engine, log *logger.Logger) {
	router.GET("/health", func(c *gin.Context) {
		HealthCheck(c)
		log.Debug("Health check performed", map[string]interface{}{
			"client_ip": c.ClientIP(),
		})
	})

	router.GET("/metrics", func(c *gin.Context) {
		// In a real application, you would expose metrics here
		// For now, we'll just return a simple response
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Metrics endpoint",
		})
		log.Debug("Metrics endpoint accessed", map[string]interface{}{
			"client_ip": c.ClientIP(),
		})
	})
}
