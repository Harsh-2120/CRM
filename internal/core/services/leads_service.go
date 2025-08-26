package services

import (
	"crm/internal/adapters/database/db"
	"crm/internal/core/repository"
	"errors"
	"fmt"

	"time"
)

// Custom errors
var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrLeadNotFound    = errors.New("lead not found")
	ErrMissingRequired = errors.New("missing required fields")
)

type LeadService interface {
	CreateLead(lead *models.Lead) (*models.Lead, error)
	UpdateLead(lead *models.Lead) (*models.Lead, error)
	GetLead(id uint) (*models.Lead, error)
	GetLeadByEmail(email string) (*models.Lead, error)
	DeleteLead(id uint) error
	GetAllLeads() ([]models.Lead, error)
}

type LeadService struct {
	queries *db.Queries
	producer KafkaProducer
}

func NewLeadService(repo repository.LeadRepository) LeadService {
	return &leadService{repo: repo}
}

// CreateLead creates a new lead in the database, with validation logic.
func (s *leadService) CreateLead(lead *models.Lead) (*models.Lead, error) {
	// Validate required fields
	if lead.FirstName == "" || lead.LastName == "" || lead.Email == "" || lead.Status == "" {
		fmt.Println("lead.FirstName: ", lead.FirstName+" lead.LastName: ", lead.LastName+" lead.Email: ", lead.Email+" lead.Status: ", lead.Status)
		return nil, ErrMissingRequired
	}

	// Validate email format
	if !isValidEmail(lead.Email) {
		return nil, ErrInvalidEmail
	}

	// Optionally: check for duplicates, or if the user is allowed to create the lead
	// Example: Check if a lead with the same email already exists
	existingLead, err := s.repo.GetByEmail(lead.Email)
	if err == nil && existingLead != nil {
		return nil, errors.New("a lead with this email already exists")
	}

	// Create lead in the repository
	fmt.Println("Creating lead in the repository")
	return s.repo.Create(lead)
}

// GetLeadByEmail retrieves a lead by its email address, with error handling.
func (s *leadService) GetLeadByEmail(email string) (*models.Lead, error) {
	lead, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if lead == nil {
		return nil, ErrLeadNotFound
	}

	return lead, nil
}

// UpdateLead updates the lead in the database, with validation logic.
func (s *leadService) UpdateLead(lead *models.Lead) (*models.Lead, error) {
	// Validate if the lead exists
	lead.UpdatedAt = time.Now()
	existingLead, err := s.repo.GetByID(lead.ID)
	if err != nil || existingLead == nil {
		return nil, ErrLeadNotFound
	}

	// Validate email if it's being changed
	if lead.Email != "" && !isValidEmail(lead.Email) {
		return nil, ErrInvalidEmail
	}

	// Perform update in repository
	return s.repo.Update(lead)
}

// GetLead retrieves a lead by its ID, with error handling.
func (s *leadService) GetLead(id uint) (*models.Lead, error) {
	lead, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if lead == nil {
		return nil, ErrLeadNotFound
	}

	return lead, nil
}

// DeleteLead removes a lead from the database, ensuring the lead exists.
func (s *leadService) DeleteLead(id uint) error {
	lead, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if lead == nil {
		return ErrLeadNotFound
	}

	return s.repo.Delete(id)
}

// GetAllLeads retrieves all leads from the database.
func (s *leadService) GetAllLeads() ([]models.Lead, error) {
	return s.repo.GetAll()
}

// Helper function to validate email addresses using a regex
// func isValidEmail(email string) bool {
// 	// Basic email validation using regex
// 	regex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
// 	return regex.MatchString(email)
// }
