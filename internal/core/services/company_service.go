package services

import (
	"crm/internal/core/domain/models"
	"crm/internal/core/repository"
	"errors"
	"strings"
	"time"
)

var (
	ErrCompanyNotFound    = errors.New("company not found")
	ErrInvalidCompanyData = errors.New("invalid company data")
)

type CompanyService interface {
	CreateCompany(company *models.Company) (*models.Company, error)
	GetCompany(id uint) (*models.Company, error)
	UpdateCompany(company *models.Company) (*models.Company, error)
	DeleteCompany(id uint) error
	ListCompanies(orgID uint, page, size uint, sortBy string, asc bool) ([]models.Company, error)
}

type companyService struct {
	repo repository.CompanyRepository
}

func NewCompanyService(repo repository.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) CreateCompany(company *models.Company) (*models.Company, error) {
	// Basic validation
	if strings.TrimSpace(company.Name) == "" || company.OrganizationID == 0 {
		return nil, ErrInvalidCompanyData
	}

	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	return s.repo.Create(company)
}

func (s *companyService) GetCompany(id uint) (*models.Company, error) {
	company, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrCompanyNotFound) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}
	return company, nil
}

func (s *companyService) UpdateCompany(company *models.Company) (*models.Company, error) {
	if company.ID == 0 || strings.TrimSpace(company.Name) == "" {
		return nil, ErrInvalidCompanyData
	}

	company.UpdatedAt = time.Now()

	updated, err := s.repo.Update(company)
	if err != nil {
		if errors.Is(err, repository.ErrCompanyNotFound) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}

	return updated, nil
}

func (s *companyService) DeleteCompany(id uint) error {
	return s.repo.Delete(id)
}

func (s *companyService) ListCompanies(orgID uint, page, size uint, sortBy string, asc bool) ([]models.Company, error) {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	return s.repo.List(orgID, page, size, sortBy, asc)
}
