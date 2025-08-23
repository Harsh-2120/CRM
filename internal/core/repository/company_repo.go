package repository

import (
	"context"
	"crm/internal/adapters/database/db"
	"errors"
)

var (
	ErrCompanyNotFound = errors.New("company not found")
)

// CompanyRepository defines the interface for company-related DB operations.
type CompanyRepository interface {
	CreateCompany(ctx context.Context, arg db.CreateCompanyParams) (db.Company, error)
	GetCompany(ctx context.Context, id int32) (db.Company, error)
	UpdateCompany(ctx context.Context, arg db.UpdateCompanyParams) (db.Company, error)
	DeleteCompany(ctx context.Context, id int32) error
	ListCompanies(ctx context.Context, organizationID int32, limit, offset int32, sortBy string, asc bool) ([]db.Company, error)
}

type companyRepository struct {
	q *db.Queries
}

func NewCompanyRepository(q *db.Queries) CompanyRepository {
	return &companyRepository{q: q}
}

// ----------------- Companies -----------------

func (r *companyRepository) CreateCompany(ctx context.Context, arg db.CreateCompanyParams) (db.Company, error) {
	return r.q.CreateCompany(ctx, arg)
}

func (r *companyRepository) GetCompany(ctx context.Context, id int32) (db.Company, error) {
	company, err := r.q.GetCompany(ctx, id)
	if err != nil {
		return db.Company{}, ErrCompanyNotFound
	}
	return company, nil
}

func (r *companyRepository) UpdateCompany(ctx context.Context, arg db.UpdateCompanyParams) (db.Company, error) {
	company, err := r.q.UpdateCompany(ctx, arg)
	if err != nil {
		return db.Company{}, ErrCompanyNotFound
	}
	return company, nil
}

func (r *companyRepository) DeleteCompany(ctx context.Context, id int32) error {
	err := r.q.DeleteCompany(ctx, id)
	if err != nil {
		return ErrCompanyNotFound
	}
	return nil
}

func (r *companyRepository) ListCompanies(ctx context.Context, organizationID int32, limit, offset int32, sortBy string, asc bool) ([]db.Company, error) {
	// sqlc doesn’t support dynamic ORDER BY out of the box
	// So, best is to write multiple queries (e.g., ListCompaniesAsc, ListCompaniesDesc)
	// or just pick one default order in your SQL.

	companies, err := r.q.ListCompanies(ctx, db.ListCompaniesParams{
		OrganizationID: organizationID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, err
	}

	return companies, nil
}
