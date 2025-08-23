package repository

import (
	"context"
	"crm/internal/adapters/database/db"
	"database/sql"
	"errors"
)

var (
	ErrOpportunityNotFound = errors.New("opportunity not found")
)

type OpportunityRepository interface {
	Create(ctx context.Context, arg db.CreateOpportunityParams) (db.Opportunity, error)
	GetByID(ctx context.Context, id int32) (db.Opportunity, error)
	Update(ctx context.Context, arg db.UpdateOpportunityParams) (db.Opportunity, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, ownerID *int32) ([]db.Opportunity, error)
}

type opportunityRepository struct {
	q *db.Queries
}

func NewOpportunityRepository(q *db.Queries) OpportunityRepository {
	return &opportunityRepository{q: q}
}

// ----------------- CRUD -----------------

func (r *opportunityRepository) Create(ctx context.Context, arg db.CreateOpportunityParams) (db.Opportunity, error) {
	return r.q.CreateOpportunity(ctx, arg)
}

func (r *opportunityRepository) GetByID(ctx context.Context, id int32) (db.Opportunity, error) {
	opp, err := r.q.GetOpportunity(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Opportunity{}, ErrOpportunityNotFound
		}
		return db.Opportunity{}, err
	}
	return opp, nil
}

func (r *opportunityRepository) Update(ctx context.Context, arg db.UpdateOpportunityParams) (db.Opportunity, error) {
	opp, err := r.q.UpdateOpportunity(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Opportunity{}, ErrOpportunityNotFound
		}
		return db.Opportunity{}, err
	}
	return opp, nil
}

func (r *opportunityRepository) Delete(ctx context.Context, id int32) error {
	err := r.q.DeleteOpportunity(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOpportunityNotFound
		}
		return err
	}
	return nil
}

func (r *opportunityRepository) List(ctx context.Context, ownerID *int32) ([]db.Opportunity, error) {
	if ownerID == nil {
		return r.q.ListOpportunities(ctx, 0)
	}
	return r.q.ListOpportunities(ctx, *ownerID)
}
