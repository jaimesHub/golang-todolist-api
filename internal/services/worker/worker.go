package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/jaimesHub/golang-todo-app/internal/services/queue"
)

// Worker processes tasks from the queue
type Worker struct {
	queue      *queue.Queue
	handlers   map[string]TaskHandler
	queueName  string
	isRunning  bool
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// TaskHandler is a function that processes a task
type TaskHandler func(task *queue.Task) error

// NewWorker creates a new worker
func NewWorker(cfg config.RedisConfig, queueName string) (*Worker, error) {
	q, err := queue.NewQueue(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Worker{
		queue:      q,
		handlers:   make(map[string]TaskHandler),
		queueName:  queueName,
		isRunning:  false,
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

// RegisterHandler registers a handler for a specific task type
func (w *Worker) RegisterHandler(taskType string, handler TaskHandler) {
	w.handlers[taskType] = handler
}

// Start starts the worker
func (w *Worker) Start() {
	if w.isRunning {
		return
	}

	w.isRunning = true

	go func() {
		for {
			select {
			case <-w.ctx.Done():
				log.Println("Worker stopped")
				return
			default:
				// Process tasks
				task, err := w.queue.Dequeue(w.queueName, 5*time.Second)
				if err != nil {
					log.Printf("Error dequeueing task: %v", err)
					time.Sleep(1 * time.Second)
					continue
				}

				if task == nil {
					// No tasks available, check scheduled tasks
					w.processScheduledTasks()
					continue
				}

				// Process task
				w.processTask(task)
			}
		}
	}()

	log.Printf("Worker started for queue: %s", w.queueName)
}

// Stop stops the worker
func (w *Worker) Stop() {
	if !w.isRunning {
		return
	}

	w.cancelFunc()
	w.isRunning = false
}

// Close closes the worker and its resources
func (w *Worker) Close() error {
	w.Stop()
	return w.queue.Close()
}

// processTask processes a task
func (w *Worker) processTask(task *queue.Task) {
	handler, exists := w.handlers[task.Type]
	if !exists {
		log.Printf("No handler registered for task type: %s", task.Type)
		return
	}

	log.Printf("Processing task: %s (type: %s)", task.ID, task.Type)

	if err := handler(task); err != nil {
		log.Printf("Error processing task %s: %v", task.ID, err)
		// In a real application, you might want to implement retry logic here
	}
}

// processScheduledTasks processes due scheduled tasks
func (w *Worker) processScheduledTasks() {
	tasks, err := w.queue.GetDueScheduledTasks()
	if err != nil {
		log.Printf("Error getting due scheduled tasks: %v", err)
		return
	}

	for _, task := range tasks {
		// Enqueue task to be processed by the worker
		_, err := w.queue.Enqueue(w.queueName, task.Type, task.Data)
		if err != nil {
			log.Printf("Error enqueueing scheduled task %s: %v", task.ID, err)
		}
	}
}
