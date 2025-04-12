package services_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jaimesHub/golang-todo-app/internal/models"
	"github.com/jaimesHub/golang-todo-app/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database interface
type MockDB struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockDB) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

// First mocks the First method
func (m *MockDB) First(dest interface{}, conds ...interface{}) error {
	args := m.Called(dest, conds)
	return args.Error(0)
}

// Where mocks the Where method
func (m *MockDB) Where(query interface{}, args ...interface{}) *MockDB {
	m.Called(query, args)
	return m
}

// Save mocks the Save method
func (m *MockDB) Save(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockDB) Delete(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	// Create a mock DB
	mockDB := new(MockDB)

	// Create a task service with the mock DB
	taskService := services.NewTaskService(mockDB)

	// Set up test data
	userID := uuid.New()
	title := "Test Task"
	description := "This is a test task"
	priority := 1
	dueDate := time.Now().Add(24 * time.Hour)

	// Set up expectations
	mockDB.On("Create", mock.AnythingOfType("*models.Task")).Return(nil)

	// Call the method being tested
	task, err := taskService.CreateTask(userID, title, description, priority, &dueDate)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, title, task.Title)
	assert.Equal(t, description, task.Description)
	assert.Equal(t, "pending", task.Status)
	assert.Equal(t, priority, task.Priority)
	assert.Equal(t, userID, task.UserID)
	mockDB.AssertExpectations(t)
}

func TestGetTaskByID(t *testing.T) {
	// Create a mock DB
	mockDB := new(MockDB)

	// Create a task service with the mock DB
	taskService := services.NewTaskService(mockDB)

	// Set up test data
	taskID := uuid.New()
	userID := uuid.New()

	// Create a task to return
	expectedTask := &models.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
		Priority:    1,
		UserID:      userID,
	}

	// Set up expectations
	mockDB.On("Where", "id = ? AND user_id = ?", []interface{}{taskID, userID}).Return(mockDB)
	mockDB.On("First", mock.AnythingOfType("*models.Task"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Copy the expected task to the destination
		task := args.Get(0).(*models.Task)
		*task = *expectedTask
	})

	// Call the method being tested
	task, err := taskService.GetTaskByID(taskID, userID)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, expectedTask.ID, task.ID)
	assert.Equal(t, expectedTask.Title, task.Title)
	assert.Equal(t, expectedTask.Description, task.Description)
	assert.Equal(t, expectedTask.Status, task.Status)
	assert.Equal(t, expectedTask.Priority, task.Priority)
	assert.Equal(t, expectedTask.UserID, task.UserID)
	mockDB.AssertExpectations(t)
}
