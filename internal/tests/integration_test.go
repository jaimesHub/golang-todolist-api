package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/handlers"
	"github.com/jaimesHub/golang-todo-app/internal/middleware"
	"github.com/jaimesHub/golang-todo-app/internal/models"
	"github.com/jaimesHub/golang-todo-app/internal/services"
	"github.com/jaimesHub/golang-todo-app/internal/services/auth"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *services.UserService, *auth.JWTService) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.Default()

	// Create a mock DB (in a real test, you would use a test database)
	// For simplicity, we'll use nil here and mock the service methods
	db := nil

	// Create services
	userService := services.NewUserService(db)
	jwtService := auth.NewJWTService(&config.JWTConfig{
		Secret:     "test-secret-key",
		Expiration: 24,
	})

	// Create handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService, jwtService)

	// Setup routes
	router.POST("/api/v1/auth/register", userHandler.Register)
	router.POST("/api/v1/auth/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		protected.GET("/users/me", userHandler.GetProfile)
	}

	return router, userService, jwtService
}

func TestUserRegistrationAndLogin(t *testing.T) {
	// Setup
	router, _, _ := setupTestRouter()

	// Test registration
	registrationPayload := `{
		"email": "test@example.com",
		"password": "password123",
		"first_name": "Test",
		"last_name": "User"
	}`

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(registrationPayload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert registration response
	assert.Equal(t, http.StatusCreated, resp.Code)

	var registrationResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &registrationResponse)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", registrationResponse["message"])

	// Test login
	loginPayload := `{
		"email": "test@example.com",
		"password": "password123"
	}`

	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert login response
	assert.Equal(t, http.StatusOK, resp.Code)

	var loginResponse map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &loginResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Login successful", loginResponse["message"])
	assert.NotEmpty(t, loginResponse["token"])
}

func TestProtectedEndpoint(t *testing.T) {
	// Setup
	router, userService, jwtService := setupTestRouter()

	// Create a test user and generate a token
	userID := uuid.New()
	token, _ := jwtService.GenerateToken(userID)

	// Mock the GetUserByID method
	// In a real test, you would use a test database or a proper mock
	// This is a simplified example
	userService.GetUserByID = func(id uuid.UUID) (*models.User, error) {
		return &models.User{
			ID:        userID,
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		}, nil
	}

	// Test accessing a protected endpoint
	req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert response
	assert.Equal(t, http.StatusOK, resp.Code)

	var profileResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &profileResponse)
	assert.NoError(t, err)

	user := profileResponse["user"].(map[string]interface{})
	assert.Equal(t, "test@example.com", user["email"])
	assert.Equal(t, "Test", user["first_name"])
	assert.Equal(t, "User", user["last_name"])
}
