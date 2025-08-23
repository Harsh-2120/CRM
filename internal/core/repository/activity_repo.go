package repository

import (
	"context"
	"crm/internal/adapters/database/db"
	"errors"
)

var (
	ErrActivityExists   = errors.New("activity with this title already exists")
	ErrActivityNotFound = errors.New("activity not found")

)

// ActivityRepository defines the methods for activity-related database operations.
type ActivityRepository interface {
	CreateActivity(ctx context.Context, arg db.CreateActivityParams) (db.Activity, error)
	GetActivity(ctx context.Context, id int32) (db.Activity, error)
	UpdateActivity(ctx context.Context, arg db.UpdateActivityParams) (db.Activity, error)
	DeleteActivity(ctx context.Context, id int32) error
	ListActivities(ctx context.Context, limit, offset int32) ([]db.Activity, error)
}

type activityRepository struct {
	q *db.Queries
}

func NewActivityRepository(q *db.Queries) ActivityRepository {
	return &activityRepository{q: q}
}

// ----------------- Activities -----------------

func (r *activityRepository) CreateActivity(ctx context.Context, arg db.CreateActivityParams) (db.Activity, error) {
	activity, err := r.q.CreateActivity(ctx, arg)
	if err != nil {
		// You can check for pg error code "23505" if needed
		return db.Activity{}, ErrActivityExists
	}
	return activity, nil
}

func (r *activityRepository) GetActivity(ctx context.Context, id int32) (db.Activity, error) {
	activity, err := r.q.GetActivity(ctx, id)
	if err != nil {
		return db.Activity{}, ErrActivityNotFound
	}
	return activity, nil
}

func (r *activityRepository) UpdateActivity(ctx context.Context, arg db.UpdateActivityParams) (db.Activity, error) {
	activity, err := r.q.UpdateActivity(ctx, arg)
	if err != nil {
		return db.Activity{}, ErrActivityNotFound
	}
	return activity, nil
}

func (r *activityRepository) DeleteActivity(ctx context.Context, id int32) error {
	err := r.q.DeleteActivity(ctx, id)
	if err != nil {
		return ErrActivityNotFound
	}
	return nil
}

func (r *activityRepository) ListActivities(ctx context.Context, limit, offset int32) ([]db.Activity, error) {
	return r.q.ListActivities(ctx, db.ListActivitiesParams{
		Limit:  limit,
		Offset: offset,
	})
}