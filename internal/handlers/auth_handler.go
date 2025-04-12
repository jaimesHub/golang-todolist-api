package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaimesHub/golang-todo-app/internal/services"
	"github.com/jaimesHub/golang-todo-app/internal/services/auth"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userService *services.UserService
	jwtService  *auth.JWTService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(userService *services.UserService, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// Register handles user registration
// Note: This is now handled by UserHandler

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.userService.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate token
	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Log activity
	if err := h.userService.LogActivity(user.ID, "login", "user", user.ID, "User logged in"); err != nil {
		// Just log the error, don't fail the login
		// logger.Error("Failed to log login activity", "error", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

// RefreshToken handles refreshing JWT tokens
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Refresh token
	newToken, err := h.jwtService.RefreshToken(input.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"token":   newToken,
	})
}
