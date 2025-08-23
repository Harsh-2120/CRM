package repository

import (
	"context"
	"crm/internal/adapters/database/db"
	"errors"
)

var (
	ErrTaskExists   = errors.New("task with this title already exists")
	ErrTaskNotFound = errors.New("task not found")
)

// ActivityRepository defines the methods for activity-related database operations.
type TaskRepository interface {
	CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error)
	GetTask(ctx context.Context, id int32) (db.Task, error)
	UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error)
	DeleteTask(ctx context.Context, id int32) error
	ListTasks(ctx context.Context, activityID int32, limit, offset int32) ([]db.Task, error)
}

type taskRepository struct {
	q *db.Queries
}

func NewTaskRepository(q *db.Queries) TaskRepository {
	return &taskRepository{q: q}
}

// ----------------- Tasks -----------------

func (r *taskRepository) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	task, err := r.q.CreateTask(ctx, arg)
	if err != nil {
		return db.Task{}, ErrTaskExists
	}
	return task, nil
}

func (r *taskRepository) GetTask(ctx context.Context, id int32) (db.Task, error) {
	task, err := r.q.GetTask(ctx, id)
	if err != nil {
		return db.Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
	task, err := r.q.UpdateTask(ctx, arg)
	if err != nil {
		return db.Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int32) error {
	err := r.q.DeleteTask(ctx, id)
	if err != nil {
		return ErrTaskNotFound
	}
	return nil
}

func (r *taskRepository) ListTasks(ctx context.Context, activityID int32, limit, offset int32) ([]db.Task, error) {
	return r.q.ListTasks(ctx, db.ListTasksParams{
		ActivityID: activityID,
		Limit:      limit,
		Offset:     offset,
	})
}
