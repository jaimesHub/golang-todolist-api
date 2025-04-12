package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jaimesHub/golang-todo-app/internal/config"
)

// Task represents a task in the queue
type Task struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}

// Queue handles task queue operations
type Queue struct {
	client *redis.Client
	ctx    context.Context
}

// NewQueue creates a new queue
func NewQueue(cfg config.RedisConfig) (*Queue, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// Ping Redis to check connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Queue{
		client: client,
		ctx:    ctx,
	}, nil
}

// Close closes the Redis connection
func (q *Queue) Close() error {
	return q.client.Close()
}

// Enqueue adds a task to the queue
func (q *Queue) Enqueue(queueName, taskType string, data map[string]interface{}) (string, error) {
	task := Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Data:      data,
		CreatedAt: time.Now(),
	}

	// Serialize task to JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to marshal task: %w", err)
	}

	// Add task to queue
	if err := q.client.RPush(q.ctx, queueName, taskJSON).Err(); err != nil {
		return "", fmt.Errorf("failed to enqueue task: %w", err)
	}

	return task.ID, nil
}

// Dequeue removes and returns a task from the queue
func (q *Queue) Dequeue(queueName string, timeout time.Duration) (*Task, error) {
	// Use BLPOP to wait for a task
	result, err := q.client.BLPop(q.ctx, timeout, queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No tasks available
		}
		return nil, fmt.Errorf("failed to dequeue task: %w", err)
	}

	// Parse task from JSON
	var task Task
	if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return &task, nil
}

// GetQueueLength returns the number of tasks in the queue
func (q *Queue) GetQueueLength(queueName string) (int64, error) {
	return q.client.LLen(q.ctx, queueName).Result()
}

// ScheduleTask schedules a task to be executed at a specific time
func (q *Queue) ScheduleTask(taskType string, data map[string]interface{}, executeAt time.Time) (string, error) {
	task := Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Data:      data,
		CreatedAt: time.Now(),
	}

	// Serialize task to JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to marshal task: %w", err)
	}

	// Add task to sorted set with score as Unix timestamp
	if err := q.client.ZAdd(q.ctx, "scheduled_tasks", &redis.Z{
		Score:  float64(executeAt.Unix()),
		Member: taskJSON,
	}).Err(); err != nil {
		return "", fmt.Errorf("failed to schedule task: %w", err)
	}

	return task.ID, nil
}

// GetDueScheduledTasks returns tasks that are due to be executed
func (q *Queue) GetDueScheduledTasks() ([]*Task, error) {
	now := time.Now().Unix()

	// Get tasks with score <= now
	results, err := q.client.ZRangeByScore(q.ctx, "scheduled_tasks", &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%d", now),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get due tasks: %w", err)
	}

	var tasks []*Task
	for _, result := range results {
		var task Task
		if err := json.Unmarshal([]byte(result), &task); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	// Remove the tasks from the sorted set
	if len(results) > 0 {
		if err := q.client.ZRemRangeByScore(q.ctx, "scheduled_tasks", "0", fmt.Sprintf("%d", now)).Err(); err != nil {
			return nil, fmt.Errorf("failed to remove due tasks: %w", err)
		}
	}

	return tasks, nil
}
