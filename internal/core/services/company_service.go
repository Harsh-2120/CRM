package services

import (
	"context"
	"crm/internal/adapters/database/db"
	"errors"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
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
	kafka   *kafka.Writer
}

func NewCompanyService(queries *db.Queries, kafkaWriter *kafka.Writer) *CompanyService {
	return &CompanyService{queries: queries, kafka: kafkaWriter}
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
	err = s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte("company_created"),
		Value: []byte(createdCompany.Name),
	})
	if err != nil {
		log.Printf("failed to write kafka message: %v", err)
	}

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
	err = s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte("company_updated"),
		Value: []byte(updatedCompany.Name),
	})
	if err != nil {
		log.Printf("failed to write kafka message: %v", err)
	}

	return &updatedCompany, nil
}

func (s *CompanyService) DeleteCompany(ctx context.Context, id int32) error {
	err := s.queries.DeleteCompany(ctx, id)
	if err != nil {
		return ErrCompanyNotFound
	}

	// Kafka Event
	err = s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte("company_deleted"),
		Value: []byte(string(rune(id))),
	})
	if err != nil {
		log.Printf("failed to write kafka message: %v", err)
	}
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
