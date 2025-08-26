package services

import (
	"context"
	"crm/internal/adapters/database/db"
	"crm/internal/adapters/kafka"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrInvalidTaskData = errors.New("invalid task data")
	ErrTaskExists      = errors.New("task with this title already exists")
)

// ActivityService defines the methods for activity and task management.
type TaskService interface {
	CreateTask(ctx context.Context, task *db.CreateTaskParams) (*db.Task, error)
	GetTask(ctx context.Context, id int32) (*db.Task, error)
	UpdateTask(ctx context.Context, params db.UpdateTaskParams) (*db.Task, error)
	DeleteTask(ctx context.Context, id int32) error
	ListTasks(ctx context.Context, pageNumber, pageSize uint) ([]db.Task, error)
}

type taskService struct {
	queries *db.Queries
	kafka   *kafka.Producer
}

func NewTaskService(queries *db.Queries, producer *kafka.Producer) TaskService {
	return &taskService{queries: queries, kafka: producer}
}

// CreateTask validates and creates a new task.
func (s *taskService) CreateTask(ctx context.Context, task *db.CreateTaskParams) (*db.Task, error) {
	if task.Title == "" || task.Status == "" || task.Priority == "" || task.ActivityID == 0 {
		return nil, ErrInvalidTaskData
	}

	validStatuses := map[string]bool{
		"Pending":     true,
		"In Progress": true,
		"Completed":   true,
	}
	if !validStatuses[task.Status] {
		return nil, errors.New("invalid task status")
	}

	validPriorities := map[string]bool{
		"Low":    true,
		"Medium": true,
		"High":   true,
	}
	if !validPriorities[task.Priority] {
		return nil, errors.New("invalid task priority")
	}

	// Validate DueDate
	if task.DueDate.Valid && task.DueDate.Time.Before(time.Now()) {
		return nil, errors.New("due date cannot be in the past")
	}
	if !task.DueDate.Valid {
		task.DueDate = sql.NullTime{Time: time.Now().Add(48 * time.Hour), Valid: true}
	}

	createdTask, err := s.queries.CreateTask(ctx, *task)
	if err != nil {
		return nil, err
	}

	// Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicTaskCreated, "task_created", map[string]interface{}{
		"id":       createdTask.ID,
		"title":    createdTask.Title,
		"status":   createdTask.Status,
		"priority": createdTask.Priority,
	})

	return &createdTask, nil
}

// GetTask retrieves a task by Id.
func (s *taskService) GetTask(ctx context.Context, id int32) (*db.Task, error) {
	task, err := s.queries.GetTask(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return &task, nil
}

// UpdateTask validates and updates an existing task.
func (s *taskService) UpdateTask(ctx context.Context, params db.UpdateTaskParams) (*db.Task, error) {
	if params.ID == 0 {
		return nil, ErrInvalidTaskData
	}

	if params.Status != "" {
		validStatuses := map[string]bool{
			"Pending":     true,
			"In Progress": true,
			"Completed":   true,
		}
		if !validStatuses[params.Status] {
			return nil, errors.New("invalid task status")
		}
	}

	if params.Priority != "" {
		validPriorities := map[string]bool{
			"Low":    true,
			"Medium": true,
			"High":   true,
		}
		if !validPriorities[params.Priority] {
			return nil, errors.New("invalid task priority")
		}
	}

	if params.DueDate.Valid && params.DueDate.Time.Before(time.Now()) {
		return nil, errors.New("due date cannot be in the past")
	}

	updatedTask, err := s.queries.UpdateTask(ctx, params)
	if err != nil {
		return nil, err
	}

	// Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicTaskUpdated, "task_updated", map[string]interface{}{
		"id":       updatedTask.ID,
		"title":    updatedTask.Title,
		"status":   updatedTask.Status,
		"priority": updatedTask.Priority,
	})

	return &updatedTask, nil
}

// DeleteTask removes a task by Id.
func (s *taskService) DeleteTask(ctx context.Context, id int32) error {
	err := s.queries.DeleteTask(ctx, id)
	if err != nil {
		return ErrTaskNotFound
	}

	// Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicTaskDeleted, "task_deleted", map[string]interface{}{
		"id": id,
	})

	return nil
}

// ListTasks retrieves tasks with pagination (basic since SQL is fixed).
func (s *taskService) ListTasks(ctx context.Context, pageNumber, pageSize uint) ([]db.Task, error) {
	if pageNumber == 0 {
		pageNumber = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	offset := (pageNumber - 1) * pageSize

	tasks, err := s.queries.ListTasks(ctx, db.ListTasksParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
