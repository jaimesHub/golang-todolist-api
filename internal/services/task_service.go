package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jaimesHub/golang-todo-app/internal/models"
	"gorm.io/gorm"
)

// TaskService handles task-related business logic
type TaskService struct {
	db *gorm.DB
}

// NewTaskService creates a new task service
func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(userID uuid.UUID, title, description string, priority int, dueDate *time.Time) (*models.Task, error) {
	// Create task
	task := &models.Task{
		Title:       title,
		Description: description,
		Status:      "pending",
		Priority:    priority,
		DueDate:     dueDate,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// GetTaskByID retrieves a task by ID
func (s *TaskService) GetTaskByID(id uuid.UUID, userID uuid.UUID) (*models.Task, error) {
	var task models.Task
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

// GetTasks retrieves tasks for a user with pagination and filtering
func (s *TaskService) GetTasks(userID uuid.UUID, status string, priority int, limit, offset int) ([]models.Task, error) {
	var tasks []models.Task

	query := s.db.Where("user_id = ?", userID)

	// Apply filters if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if priority >= 0 {
		query = query.Where("priority = ?", priority)
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	// Order by priority (high to low) and created_at (newest first)
	query = query.Order("priority DESC, created_at DESC")

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(id uuid.UUID, userID uuid.UUID, title, description, status string, priority int, dueDate *time.Time) (*models.Task, error) {
	// Get task
	task, err := s.GetTaskByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if title != "" {
		task.Title = title
	}

	if description != "" {
		task.Description = description
	}

	if status != "" {
		task.Status = status
	}

	if priority >= 0 {
		task.Priority = priority
	}

	if dueDate != nil {
		task.DueDate = dueDate
	}

	task.UpdatedAt = time.Now()

	if err := s.db.Save(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(id uuid.UUID, userID uuid.UUID) error {
	// Check if task exists and belongs to user
	task, err := s.GetTaskByID(id, userID)
	if err != nil {
		return err
	}

	// Delete task (soft delete with GORM)
	if err := s.db.Delete(task).Error; err != nil {
		return err
	}

	return nil
}

// CountTasks counts tasks for a user with optional status filter
func (s *TaskService) CountTasks(userID uuid.UUID, status string) (int64, error) {
	var count int64

	query := s.db.Model(&models.Task{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
