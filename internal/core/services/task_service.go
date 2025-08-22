package services

import (
	"crm/internal/core/domain/models"
	"crm/internal/core/repository"
	"errors"
	"regexp"
	"time"
)

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrInvalidTaskData = errors.New("invalid task data")
	ErrTaskExists      = errors.New("task with this title already exists")
)

// ActivityService defines the methods for activity and task management.
type TaskService interface {
	CreateTask(task *models.Task) (*models.Task, error)
	GetTask(id uint) (*models.Task, error)
	UpdateTask(task *models.Task) (*models.Task, error)
	DeleteTask(id uint) error
	ListTasks(pageNumber uint, pageSize uint, sortBy string, ascending bool, activityID uint) ([]models.Task, error)
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// CreateTask validates and creates a new task.
func (s *taskService) CreateTask(task *models.Task) (*models.Task, error) {
	// Validate required fields
	if task.Title == "" || task.Status == "" || task.Priority == "" || task.ActivityID == 0 {
		return nil, ErrInvalidTaskData
	}

	// Validate Status and Priority against predefined sets
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
	if !task.DueDate.IsZero() && task.DueDate.Before(time.Now()) {
		return nil, errors.New("due date cannot be in the past")
	}

	// Set timestamps
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	// Attempt to create the task
	createdTask, err := s.repo.CreateTask(task)
	if err != nil {
		if errors.Is(err, repository.ErrTaskExists) {
			return nil, ErrTaskExists
		}
		return nil, err
	}

	return createdTask, nil
}

// GetTask retrieves a task by Id.
func (s *taskService) GetTask(id uint) (*models.Task, error) {
	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}

// UpdateTask validates and updates an existing task.
func (s *taskService) UpdateTask(task *models.Task) (*models.Task, error) {
	// Validate task Id
	if task.Id == 0 {
		return nil, ErrInvalidTaskData
	}

	// Validate Status and Priority if provided
	if task.Status != "" {
		validStatuses := map[string]bool{
			"Pending":     true,
			"In Progress": true,
			"Completed":   true,
		}
		if !validStatuses[task.Status] {
			return nil, errors.New("invalid task status")
		}
	}

	if task.Priority != "" {
		validPriorities := map[string]bool{
			"Low":    true,
			"Medium": true,
			"High":   true,
		}
		if !validPriorities[task.Priority] {
			return nil, errors.New("invalid task priority")
		}
	}

	// Validate DueDate if provided
	if !task.DueDate.IsZero() && task.DueDate.Before(time.Now()) {
		return nil, errors.New("due date cannot be in the past")
	}

	// Set the UpdatedAt timestamp
	task.UpdatedAt = time.Now()

	// Update the task
	updatedTask, err := s.repo.UpdateTask(task)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}
		if errors.Is(err, repository.ErrTaskExists) {
			return nil, ErrTaskExists
		}
		return nil, err
	}

	return updatedTask, nil
}

// DeleteTask removes a task by Id.
func (s *taskService) DeleteTask(id uint) error {
	// Check if the task exists
	_, err := s.repo.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	// Delete the task
	if err := s.repo.DeleteTask(id); err != nil {
		return err
	}
	return nil
}

// ListTasks retrieves tasks with pagination, sorting, and optional filtering by activity.
func (s *taskService) ListTasks(pageNumber uint, pageSize uint, sortBy string, ascending bool, activityID uint) ([]models.Task, error) {
	// Validate pagination parameters
	if pageNumber == 0 {
		pageNumber = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	// Validate sortBy field
	validSortFields := map[string]bool{
		"title":      true,
		"due_date":   true,
		"created_at": true,
		"updated_at": true,
		"status":     true,
		"priority":   true,
	}
	if sortBy != "" && !validSortFields[sortBy] {
		return nil, errors.New("invalid sort field")
	}

	tasks, err := s.repo.ListTasks(pageNumber, pageSize, sortBy, ascending, activityID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Helper function to validate email format using regex (if needed for tasks).
func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
