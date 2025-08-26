package services

import (
	"context"
	"crm/internal/adapters/database/db"
	"crm/internal/adapters/kafka"
	"errors"
	"strings"
)

var (
	ErrCompanyNotFound    = errors.New("company not found")
	ErrInvalidCompanyData = errors.New("invalid company data")
)

type CompanyServiceInterface interface {
	CreateCompany(ctx context.Context, company db.CreateCompanyParams) (*db.Company, error)
	GetCompany(ctx context.Context, id int32) (*db.Company, error)
	UpdateCompany(ctx context.Context, company db.UpdateCompanyParams) (*db.Company, error)
	DeleteCompany(ctx context.Context, id int32) error
	ListCompanies(ctx context.Context, orgID int32, page, size int32) ([]db.Company, error)
}
type CompanyService struct {
	queries *db.Queries
	kafka   *kafka.Producer
}

func NewCompanyService(queries *db.Queries, producer *kafka.Producer) *CompanyService {
	return &CompanyService{queries: queries, kafka: producer}
}

func (s *CompanyService) CreateCompany(ctx context.Context, company db.CreateCompanyParams) (*db.Company, error) {
	if strings.TrimSpace(company.Name) == "" || company.OrganizationID == 0 {
		return nil, ErrInvalidCompanyData
	}

	createdCompany, err := s.queries.CreateCompany(ctx, company)
	if err != nil {
		return nil, err
	}

	// Publish Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicCompanyCreated, "company_created", map[string]interface{}{
		"id":   createdCompany.ID,
		"name": createdCompany.Name,
	})

	return &createdCompany, nil
}

func (s *CompanyService) GetCompany(ctx context.Context, id int32) (*db.Company, error) {
	company, err := s.queries.GetCompany(ctx, id)
	if err != nil {
		return nil, ErrCompanyNotFound
	}
	return &company, nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, company db.UpdateCompanyParams) (*db.Company, error) {
	if company.ID == 0 || strings.TrimSpace(company.Name) == "" {
		return nil, ErrInvalidCompanyData
	}

	updatedCompany, err := s.queries.UpdateCompany(ctx, company)
	if err != nil {
		return nil, ErrCompanyNotFound
	}

	// Kafka Event
	_ = s.kafka.Publish(ctx, kafka.TopicCompanyUpdated, "company_updated", map[string]interface{}{
		"id":   updatedCompany.ID,
		"name": updatedCompany.Name,
	})

	return &updatedCompany, nil
}

func (s *CompanyService) DeleteCompany(ctx context.Context, id int32) error {
	err := s.queries.DeleteCompany(ctx, id)
	if err != nil {
		return ErrCompanyNotFound
	}

	// Kafka Event
	_ = s.kafka.Publish(ctx, kafka.TopicCompanyDeleted, "company_deleted", map[string]interface{}{
		"id": id,
	})

	return nil
}

func (s *CompanyService) ListCompanies(ctx context.Context, orgID int32, page, size uint) ([]db.Company, error) {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	offset := (page - 1) * size
	companies, err := s.queries.ListCompanies(ctx, db.ListCompaniesParams{
		OrganizationID: orgID,
		Limit:          int32(size),
		Offset:         int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return companies, nil
}
