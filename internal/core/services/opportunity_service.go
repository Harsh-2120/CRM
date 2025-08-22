package services

import (
	"crm/internal/core/domain/models"
	"crm/internal/core/repository"
	"errors"

	"time"

	"gorm.io/gorm"
)

var (
	ErrOpportunityNotFound    = errors.New("opportunity not found")
	ErrInvalidOpportunityData = errors.New("invalid opportunity data")
)

type OpportunityService interface {
	CreateOpportunity(opportunity *models.Opportunity) (*models.Opportunity, error)
	GetOpportunity(id uint) (*models.Opportunity, error)
	UpdateOpportunity(opportunity *models.Opportunity) (*models.Opportunity, error)
	DeleteOpportunity(id uint) error
	ListOpportunities(ownerID uint) ([]models.Opportunity, error)
}

type opportunityService struct {
	repo repository.OpportunityRepository
}

func NewOpportunityService(repo repository.OpportunityRepository) OpportunityService {
	return &opportunityService{repo}
}

func (s *opportunityService) CreateOpportunity(opportunity *models.Opportunity) (*models.Opportunity, error) {
	// Validate required fields
	print("in service call\n")
	if opportunity.CloseDate.IsZero() {
		panic("failed to validate closedate")
	}
	print("validated date")

	if opportunity.Name == "" || opportunity.Stage == "" || opportunity.Amount == 0 || opportunity.CloseDate.IsZero() || opportunity.OwnerID == 0 {
		return nil, ErrInvalidOpportunityData
	}

	// Additional validations
	if opportunity.Probability != 0 && (opportunity.Probability < 0 || opportunity.Probability > 100) {
		return nil, errors.New("probability must be between 0 and 100")
	}

	// Set timestamps
	now := time.Now()
	opportunity.CreatedAt = now
	opportunity.UpdatedAt = now

	return s.repo.Create(opportunity)
}

func (s *opportunityService) GetOpportunity(id uint) (*models.Opportunity, error) {
	opportunity, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOpportunityNotFound
		}
		return nil, err
	}
	return opportunity, nil
}

func (s *opportunityService) UpdateOpportunity(opportunity *models.Opportunity) (*models.Opportunity, error) {
	if opportunity.Id == 0 {
		return nil, ErrInvalidOpportunityData
	}

	// Set the UpdatedAt field
	opportunity.UpdatedAt = time.Now()

	// Perform the update, omitting zero values
	err := s.repo.UpdateSelective(opportunity)
	if err != nil {
		return nil, err
	}

	// Retrieve the updated opportunity
	updatedOpportunity, err := s.repo.GetByID(opportunity.Id)
	if err != nil {
		return nil, err
	}

	return updatedOpportunity, nil
}

func (s *opportunityService) DeleteOpportunity(id uint) error {
	// Check if the opportunity exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOpportunityNotFound
		}
		return err
	}
	return s.repo.Delete(id)
}

func (s *opportunityService) ListOpportunities(ownerID uint) ([]models.Opportunity, error) {
	return s.repo.List(ownerID)
}
