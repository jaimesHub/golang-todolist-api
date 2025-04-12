package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jaimesHub/golang-todo-app/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService handles user-related business logic
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(email, password, firstName, lastName string) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	result := s.db.Where("email = ?", email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("user with this email already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// Log activity
	activity := &models.Activity{
		UserID:    user.ID,
		Action:    "register",
		Entity:    "user",
		EntityID:  user.ID,
		Details:   "User registration",
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(activity).Error; err != nil {
		// Just log the error, don't fail the user creation
		// In a real app, you might want to use a proper logger
		// logger.Error("Failed to log user registration activity", "error", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user's profile
func (s *UserService) UpdateUser(id uuid.UUID, firstName, lastName string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields
	user.FirstName = firstName
	user.LastName = lastName
	user.UpdatedAt = time.Now()

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	// Log activity
	activity := &models.Activity{
		UserID:    user.ID,
		Action:    "update",
		Entity:    "user",
		EntityID:  user.ID,
		Details:   "User profile updated",
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(activity).Error; err != nil {
		// Just log the error, don't fail the user update
		// logger.Error("Failed to log user update activity", "error", err)
	}

	return &user, nil
}

// GetUserActivities retrieves a user's activities
func (s *UserService) GetUserActivities(userID uuid.UUID, limit, offset int) ([]models.Activity, error) {
	var activities []models.Activity

	query := s.db.Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

// LogActivity logs a user activity
func (s *UserService) LogActivity(userID uuid.UUID, action, entity string, entityID uuid.UUID, details string) error {
	activity := &models.Activity{
		UserID:    userID,
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		Details:   details,
		CreatedAt: time.Now(),
	}

	return s.db.Create(activity).Error
}
