package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jaimesHub/golang-todo-app/internal/services"
)

// TaskHandler handles task-related requests
type TaskHandler struct {
	taskService *services.TaskService
	userService *services.UserService
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskService *services.TaskService, userService *services.UserService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		userService: userService,
	}
}

// Create handles creating a new task
func (h *TaskHandler) Create(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description"`
		Priority    int        `json:"priority"`
		DueDate     *time.Time `json:"due_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.CreateTask(
		userID.(uuid.UUID),
		input.Title,
		input.Description,
		input.Priority,
		input.DueDate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log activity
	h.userService.LogActivity(
		userID.(uuid.UUID),
		"create",
		"task",
		task.ID,
		"Task created: "+task.Title,
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task":    task,
	})
}

// List handles listing all tasks
func (h *TaskHandler) List(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse query parameters
	status := c.Query("status")

	priority := -1 // Default value to indicate no filter
	if priorityParam := c.Query("priority"); priorityParam != "" {
		if parsedPriority, err := strconv.Atoi(priorityParam); err == nil {
			priority = parsedPriority
		}
	}

	limit := 10 // Default limit
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetParam := c.Query("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	tasks, err := h.taskService.GetTasks(userID.(uuid.UUID), status, priority, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get total count for pagination
	totalCount, err := h.taskService.CountTasks(userID.(uuid.UUID), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"pagination": gin.H{
			"total":  totalCount,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetByID handles getting a task by ID
func (h *TaskHandler) GetByID(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse task ID from URL
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.taskService.GetTaskByID(taskID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": task,
	})
}

// Update handles updating a task
func (h *TaskHandler) Update(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse task ID from URL
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var input struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		Priority    int        `json:"priority"`
		DueDate     *time.Time `json:"due_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.UpdateTask(
		taskID,
		userID.(uuid.UUID),
		input.Title,
		input.Description,
		input.Status,
		input.Priority,
		input.DueDate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log activity
	h.userService.LogActivity(
		userID.(uuid.UUID),
		"update",
		"task",
		task.ID,
		"Task updated: "+task.Title,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"task":    task,
	})
}

// Delete handles deleting a task
func (h *TaskHandler) Delete(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse task ID from URL
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get task before deletion for activity logging
	task, err := h.taskService.GetTaskByID(taskID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := h.taskService.DeleteTask(taskID, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log activity
	h.userService.LogActivity(
		userID.(uuid.UUID),
		"delete",
		"task",
		taskID,
		"Task deleted: "+task.Title,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}
