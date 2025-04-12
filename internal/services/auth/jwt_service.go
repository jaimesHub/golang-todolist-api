package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jaimesHub/golang-todo-app/internal/config"
)

// JWTService handles JWT token generation and validation
type JWTService struct {
	config *config.JWTConfig
}

// TokenClaims represents the JWT claims
type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

// NewJWTService creates a new JWT service
func NewJWTService(config *config.JWTConfig) *JWTService {
	return &JWTService{config: config}
}

// GenerateToken generates a new JWT token for a user
func (s *JWTService) GenerateToken(userID uuid.UUID) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(time.Hour * time.Duration(s.config.Expiration))

	// Create claims
	claims := &TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "todo-app",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken refreshes a JWT token
func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	// Validate token
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same user ID
	return s.GenerateToken(claims.UserID)
}
