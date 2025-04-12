package services_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/services/auth"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// Create a JWT service with test configuration
	jwtConfig := &config.JWTConfig{
		Secret:     "test-secret-key",
		Expiration: 24,
	}
	jwtService := auth.NewJWTService(jwtConfig)

	// Generate a token for a test user
	userID := uuid.New()
	token, err := jwtService.GenerateToken(userID)

	// Assert token generation
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims, err := jwtService.ValidateToken(token)

	// Assert token validation
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)

	// Check that expiration is set correctly (approximately)
	expectedExpiration := time.Now().Add(time.Hour * 24).Unix()
	assert.InDelta(t, expectedExpiration, claims.ExpiresAt, 10) // Allow 10 seconds delta
}

func TestInvalidToken(t *testing.T) {
	// Create a JWT service with test configuration
	jwtConfig := &config.JWTConfig{
		Secret:     "test-secret-key",
		Expiration: 24,
	}
	jwtService := auth.NewJWTService(jwtConfig)

	// Test with invalid token
	_, err := jwtService.ValidateToken("invalid-token")
	assert.Error(t, err)

	// Test with expired token
	// Create a token that's already expired
	userID := uuid.New()
	claims := &auth.TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
			IssuedAt:  time.Now().Add(-time.Hour * 2).Unix(),
			Issuer:    "todo-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, _ := token.SignedString([]byte(jwtConfig.Secret))

	_, err = jwtService.ValidateToken(expiredToken)
	assert.Error(t, err)
}

func TestRefreshToken(t *testing.T) {
	// Create a JWT service with test configuration
	jwtConfig := &config.JWTConfig{
		Secret:     "test-secret-key",
		Expiration: 24,
	}
	jwtService := auth.NewJWTService(jwtConfig)

	// Generate a token for a test user
	userID := uuid.New()
	token, err := jwtService.GenerateToken(userID)
	assert.NoError(t, err)

	// Refresh the token
	newToken, err := jwtService.RefreshToken(token)

	// Assert token refresh
	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)
	assert.NotEqual(t, token, newToken) // New token should be different

	// Validate the new token
	claims, err := jwtService.ValidateToken(newToken)

	// Assert new token validation
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
}
