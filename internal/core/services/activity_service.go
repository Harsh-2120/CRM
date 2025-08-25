package services

import (
	"context"
	"crm/internal/adapters/database/db"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	ErrActivityNotFound    = errors.New("activity not found")
	ErrInvalidActivityData = errors.New("invalid activity data")
)

type ActivityService interface {
	CreateActivity(ctx context.Context, activity *db.CreateActivityParams) (*db.Activity, error)
	GetActivity(ctx context.Context, id int32) (*db.Activity, error)
	UpdateActivity(ctx context.Context, params db.UpdateActivityParams) (*db.Activity, error)
	DeleteActivity(ctx context.Context, id int32) error
	ListActivities(ctx context.Context, pageNumber, pageSize uint) ([]db.Activity, error)
}

type activityService struct {
	queries *db.Queries
	kafka   *kafka.Writer
}

func NewActivityService(queries *db.Queries, kafkaWriter *kafka.Writer) ActivityService {
	return &activityService{queries: queries, kafka: kafkaWriter}
}

// CreateActivity validates and creates a new activity.
func (s *activityService) CreateActivity(ctx context.Context, activity *db.CreateActivityParams) (*db.Activity, error) {
	if activity.Title == "" || activity.Type == "" || activity.Status == "" || activity.ContactID == 0 {
		return nil, ErrInvalidActivityData
	}

	// If due_date is not set, assign default
	if !activity.DueDate.Valid {
		activity.DueDate = sql.NullTime{
			Time:  time.Now().Add(24 * time.Hour),
			Valid: true,
		}
	}

	createdActivity, err := s.queries.CreateActivity(ctx, *activity)
	if err != nil {
		return nil, err
	}

	// Publish Kafka Event
	_ = s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte("activity_created"),
		Value: []byte(createdActivity.Title),
	})

	return &createdActivity, nil
}

// GetActivity retrieves an activity by Id.
func (s *activityService) GetActivity(ctx context.Context, id int32) (*db.Activity, error) {
	activity, err := s.queries.GetActivity(ctx, id)
	if err != nil {
		return nil, ErrActivityNotFound
	}
	return &activity, nil
}

// UpdateActivity validates and updates an existing activity.
func (s *activityService) UpdateActivity(ctx context.Context, params db.UpdateActivityParams) (*db.Activity, error) {
	if params.ID == 0 {
		return nil, ErrInvalidActivityData
	}

	if params.Status != "" {
		validStatuses := map[string]bool{
			"Pending":   true,
			"Completed": true,
			"Canceled":  true,
		}
		if !validStatuses[params.Status] {
			return nil, errors.New("invalid activity status")
		}
	}

	// If due date is provided, validate it
	if params.DueDate.Valid {
		if params.DueDate.Time.Before(time.Now()) {
			return nil, errors.New("due date cannot be in the past")
		}
	}

	updatedActivity, err := s.queries.UpdateActivity(ctx, params)
	if err != nil {
		return nil, err
	}

	// Publish Kafka Event
	_ = s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte("activity_updated"),
		Value: []byte(updatedActivity.Title),
	})

	return &updatedActivity, nil
}

// DeleteActivity removes an activity by Id.
func (s *activityService) DeleteActivity(ctx context.Context, id int32) error {
	err := s.queries.DeleteActivity(ctx, id)
	if err != nil {
		return ErrActivityNotFound
	}

	// Publish Kafka Event (convert id properly to string)
	_ = s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte("activity_deleted"),
		Value: []byte(fmt.Sprintf("%d", id)),
	})
	return nil
}

// ListActivities retrieves activities with pagination.
func (s *activityService) ListActivities(ctx context.Context, pageNumber, pageSize uint) ([]db.Activity, error) {
	if pageNumber == 0 {
		pageNumber = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	offset := (pageNumber - 1) * pageSize

	activities, err := s.queries.ListActivities(ctx, db.ListActivitiesParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return activities, nil
}
