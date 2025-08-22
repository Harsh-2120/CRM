package services

import (
	"context"
	"crm/internal/core/domain/models"
	"crm/internal/core/repository"
	"errors"
	"log"
	"time"
)

var (
	ErrActivityNotFound    = errors.New("activity not found")
	ErrInvalidActivityData = errors.New("invalid activity data")
	ErrActivityExists      = errors.New("activity with this title already exists")
)

// ActivityService defines the methods for activity and task management.
type ActivityService interface {
	CreateActivity(ctx context.Context, activity *models.Activity) (*models.Activity, error)
	GetActivity(id uint) (*models.Activity, error)
	UpdateActivity(activity *models.Activity) (*models.Activity, error)
	DeleteActivity(id uint) error
	ListActivities(pageNumber, pageSize uint, sortBy string, ascending bool, contactID uint) ([]models.Activity, error)
	GetActivityByID(id uint) (*models.Activity, error)
}

type activityService struct {
	repo repository.ActivityRepository
}

func NewActivityService(repo repository.ActivityRepository) ActivityService {
	return &activityService{repo: repo}
}

// GetActivityByID implements ActivityService.
func (s *activityService) GetActivityByID(id uint) (*models.Activity, error) {
	activity, err := s.repo.GetActivityByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrActivityNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}
	return activity, nil
}

// CreateActivity validates and creates a new activity.
func (s *activityService) CreateActivity(ctx context.Context, activity *models.Activity) (*models.Activity, error) {
	// Validate required fields
	if activity.Title == "" || activity.Type == "" || activity.Status == "" || activity.ContactID == 0 {
		return nil, ErrInvalidActivityData
	}

	// Validate ActivityStatus using constants or enums
	validStatuses := map[string]bool{
		"Pending":    true,
		"InProgress": true,
		"Completed":  true,
		"Canceled":   true,
		"Scheduled":  true, // Include if 'Scheduled' is a valid status
	}

	if !validStatuses[activity.Status] {
		return nil, errors.New("invalid activity status")
	}

	// Set timestamps
	now := time.Now()
	activity.CreatedAt = now
	activity.UpdatedAt = now

	// Attempt to create the activity
	createdActivity, err := s.repo.CreateActivity(ctx, activity)
	if err != nil {
		if errors.Is(err, repository.ErrActivityExists) {
			return nil, ErrActivityExists
		}
		return nil, err
	}

	return createdActivity, nil
}

// GetActivity retrieves an activity by Id.
func (s *activityService) GetActivity(id uint) (*models.Activity, error) {
	activity, err := s.repo.GetActivityByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrActivityNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}
	return activity, nil
}

// UpdateActivity validates and updates an existing activity.
func (s *activityService) UpdateActivity(activity *models.Activity) (*models.Activity, error) {
	// Validate activity Id
	if activity.Id == 0 {
		return nil, ErrInvalidActivityData
	}

	// Validate Type and Status if provided
	if activity.Type != "" {
		validTypes := map[string]bool{
			"Call":    true,
			"Meeting": true,
			"Email":   true,
		}
		if !validTypes[activity.Type] {
			return nil, errors.New("invalid activity type")
		}
	}

	if activity.Status != "" {
		validStatuses := map[string]bool{
			"Pending":   true,
			"Completed": true,
			"Canceled":  true,
		}
		if !validStatuses[activity.Status] {
			return nil, errors.New("invalid activity status")
		}
	}

	// Validate DueDate if provided
	if !activity.DueDate.IsZero() && activity.DueDate.Before(time.Now()) {
		return nil, errors.New("due date cannot be in the past")
	}

	// Set the UpdatedAt timestamp
	activity.UpdatedAt = time.Now()

	// Update the activity
	log.Printf("serice is clear on new value %v", activity)

	updatedActivity, err := s.repo.UpdateActivity(activity)
	if err != nil {
		if errors.Is(err, repository.ErrActivityNotFound) {
			return nil, ErrActivityNotFound
		}
		if errors.Is(err, repository.ErrActivityExists) {
			return nil, ErrActivityExists
		}
		return nil, err
	}

	return updatedActivity, nil
}

// DeleteActivity removes an activity by Id.
func (s *activityService) DeleteActivity(id uint) error {
	// Check if the activity exists
	_, err := s.repo.GetActivityByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrActivityNotFound) {
			return ErrActivityNotFound
		}
		return err
	}

	// Delete the activity
	if err := s.repo.DeleteActivity(id); err != nil {
		return err
	}
	return nil
}

// ListActivities retrieves activities with pagination, sorting, and optional filtering by contact.
func (s *activityService) ListActivities(pageNumber uint, pageSize uint, sortBy string, ascending bool, contactID uint) ([]models.Activity, error) {
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
		"type":       true,
		"status":     true,
	}
	if sortBy != "" && !validSortFields[sortBy] {
		return nil, errors.New("invalid sort field")
	}

	activities, err := s.repo.ListActivities(pageNumber, pageSize, sortBy, ascending, contactID)
	if err != nil {
		return nil, err
	}
	return activities, nil
}
