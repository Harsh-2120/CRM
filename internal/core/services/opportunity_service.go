package services

import (
	"context"
	"crm/internal/adapters/database/db"
	"crm/internal/adapters/kafka"
	"errors"
	"strings"
)

var (
	ErrOpportunityNotFound    = errors.New("opportunity not found")
	ErrInvalidOpportunityData = errors.New("invalid opportunity data")
)

type OpportunityServiceInterface interface {
	CreateOpportunity(ctx context.Context, opportunity db.CreateOpportunityParams) (*db.Opportunity, error)
	GetOpportunity(ctx context.Context, id int32) (*db.Opportunity, error)
	UpdateOpportunity(ctx context.Context, opportunity db.UpdateOpportunityParams) (*db.Opportunity, error)
	DeleteOpportunity(ctx context.Context, id int32) error
	ListOpportunities(ctx context.Context, ownerID int32) ([]db.Opportunity, error)
}

type OpportunityService struct {
	queries *db.Queries
	kafka   *kafka.Producer
}

func NewOpportunityService(queries *db.Queries, producer *kafka.Producer) *OpportunityService {
	return &OpportunityService{queries: queries, kafka: producer}
}

func (s *OpportunityService) CreateOpportunity(ctx context.Context, opportunity db.CreateOpportunityParams) (*db.Opportunity, error) {
	// Validate name
	if !opportunity.Name.Valid || strings.TrimSpace(opportunity.Name.String) == "" {
		return nil, ErrInvalidOpportunityData
	}

	// Validate stage
	if !opportunity.Stage.Valid || strings.TrimSpace(opportunity.Stage.String) == "" {
		return nil, ErrInvalidOpportunityData
	}

	// Validate other required fields
	if opportunity.Amount <= 0 {
		return nil, ErrInvalidOpportunityData
	}

	// Probability check (if not NULL)
	if opportunity.Probability < 0 || opportunity.Probability > 100 {
		return nil, errors.New("probability must be between 0 and 100")
	}

	createdOpportunity, err := s.queries.CreateOpportunity(ctx, opportunity)
	if err != nil {
		return nil, err
	}

	// Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicOpportunityCreated, "opportunity_created", map[string]interface{}{
		"id":     createdOpportunity.ID,
		"name":   createdOpportunity.Name.String,
		"amount": createdOpportunity.Amount,
		"stage":  createdOpportunity.Stage.String,
	})

	return &createdOpportunity, nil
}

// GetOpportunity retrieves an opportunity by ID
func (s *OpportunityService) GetOpportunity(ctx context.Context, id int32) (*db.Opportunity, error) {
	opportunity, err := s.queries.GetOpportunity(ctx, id)
	if err != nil {
		return nil, ErrOpportunityNotFound
	}
	return &opportunity, nil
}

// UpdateOpportunity validates and updates an opportunity
func (s *OpportunityService) UpdateOpportunity(ctx context.Context, opportunity db.UpdateOpportunityParams) (*db.Opportunity, error) {
	if opportunity.ID == 0 {
		return nil, ErrInvalidOpportunityData
	}

	updatedOpportunity, err := s.queries.UpdateOpportunity(ctx, opportunity)
	if err != nil {
		return nil, ErrOpportunityNotFound
	}

	// Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicOpportunityUpdated, "opportunity_updated", map[string]interface{}{
		"id":     updatedOpportunity.ID,
		"name":   updatedOpportunity.Name,
		"amount": updatedOpportunity.Amount,
		"stage":  updatedOpportunity.Stage,
	})

	return &updatedOpportunity, nil
}

// DeleteOpportunity removes an opportunity by ID
func (s *OpportunityService) DeleteOpportunity(ctx context.Context, id int32) error {
	err := s.queries.DeleteOpportunity(ctx, id)
	if err != nil {
		return ErrOpportunityNotFound
	}

	// Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicOpportunityDeleted, "opportunity_deleted", map[string]interface{}{
		"id": id,
	})

	return nil
}

// ListOpportunities returns opportunities for a given owner
func (s *OpportunityService) ListOpportunities(ctx context.Context, ownerID int32) ([]db.Opportunity, error) {
	opportunities, err := s.queries.ListOpportunities(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return opportunities, nil
}
