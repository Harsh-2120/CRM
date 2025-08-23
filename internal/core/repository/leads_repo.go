package repository

import (
	"context"
	"crm/internal/adapters/database/db"
	"database/sql"
	"errors"
)

var (
	ErrLeadNotFound = errors.New("lead not found")
)

// LeadRepository defines the contract for lead CRUD operations.
type LeadRepository interface {
	Create(ctx context.Context, arg db.CreateLeadParams) (db.Lead, error)
	GetByID(ctx context.Context, id int32) (db.Lead, error)
	GetByEmail(ctx context.Context, email string) (db.Lead, error)
	Update(ctx context.Context, arg db.UpdateLeadParams) (db.Lead, error)
	Delete(ctx context.Context, id int32) error
	GetAll(ctx context.Context) ([]db.Lead, error)
}

type leadRepository struct {
	q *db.Queries
}

func NewLeadRepository(q *db.Queries) LeadRepository {
	return &leadRepository{q: q}
}

// ----------------- CRUD -----------------

func (r *leadRepository) Create(ctx context.Context, arg db.CreateLeadParams) (db.Lead, error) {
	return r.q.CreateLead(ctx, arg)
}

func (r *leadRepository) GetByID(ctx context.Context, id int32) (db.Lead, error) {
	lead, err := r.q.GetLead(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Lead{}, ErrLeadNotFound
		}
		return db.Lead{}, err
	}
	return lead, nil
}

func (r *leadRepository) GetByEmail(ctx context.Context, email string) (db.Lead, error) {
	lead, err := r.q.GetLeadByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Lead{}, ErrLeadNotFound
		}
		return db.Lead{}, err
	}
	return lead, nil
}

func (r *leadRepository) Update(ctx context.Context, arg db.UpdateLeadParams) (db.Lead, error) {
	lead, err := r.q.UpdateLead(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Lead{}, ErrLeadNotFound
		}
		return db.Lead{}, err
	}
	return lead, nil
}

func (r *leadRepository) Delete(ctx context.Context, id int32) error {
	err := r.q.DeleteLead(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLeadNotFound
		}
		return err
	}
	return nil
}

func (r *leadRepository) GetAll(ctx context.Context) ([]db.Lead, error) {
	return r.q.ListLeads(ctx)
}
