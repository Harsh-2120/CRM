package services

import (
	"context"
	"crm/internal/adapters/database/db"
	"crm/internal/adapters/kafka"
	"errors"
	"strings"
)

var (
	ErrLeadNotFound    = errors.New("lead not found")
	ErrInvalidLeadData = errors.New("invalid lead data")
	ErrInvalidEmail    = errors.New("invalid email format")
)

type LeadServiceInterface interface {
	CreateLead(ctx context.Context, lead db.CreateLeadParams) (*db.Lead, error)
	GetLead(ctx context.Context, id int32) (*db.Lead, error)
	UpdateLead(ctx context.Context, lead db.UpdateLeadParams) (*db.Lead, error)
	DeleteLead(ctx context.Context, id int32) error
	GetAllLeads(ctx context.Context, pageNumber, pageSize int32) ([]db.Lead, error)
	GetLeadByEmail(ctx context.Context, email string) (*db.Lead, error)
}

type LeadService struct {
	queries *db.Queries
	kafka   *kafka.Producer
}

func NewLeadService(queries *db.Queries, producer *kafka.Producer) *LeadService {
	return &LeadService{queries: queries, kafka: producer}
}

// CreateLead validates and creates a new lead.
func (s *LeadService) CreateLead(ctx context.Context, lead db.CreateLeadParams) (*db.Lead, error) {
	// Required fields
	if strings.TrimSpace(lead.FirstName) == "" ||
		strings.TrimSpace(lead.LastName) == "" ||
		strings.TrimSpace(lead.Email) == "" ||
		strings.TrimSpace(lead.Status) == "" {
		return nil, ErrInvalidLeadData
	}

	// Email format
	if !isValidEmail(lead.Email) {
		return nil, ErrInvalidEmail
	}

	created, err := s.queries.CreateLead(ctx, lead)
	if err != nil {
		return nil, err
	}

	// Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicLeadCreated, "lead_created", map[string]interface{}{
		"id":     created.ID,
		"email":  created.Email,
		"status": created.Status,
	})

	return &created, nil
}

// GetLead retrieves a lead by ID.
func (s *LeadService) GetLead(ctx context.Context, id int32) (*db.Lead, error) {
	lead, err := s.queries.GetLeadById(ctx, id)
	if err != nil {
		return nil, ErrLeadNotFound
	}
	return &lead, nil
}

// UpdateLead validates and updates an existing lead.
// NOTE: sqlc UpdateLead sets status and assigned_to; ensure ID and (optionally) status sanity.
func (s *LeadService) UpdateLead(ctx context.Context, lead db.UpdateLeadParams) (*db.Lead, error) {
	if lead.ID == 0 {
		return nil, ErrInvalidLeadData
	}
	if strings.TrimSpace(lead.Status) == "" {
		return nil, ErrInvalidLeadData
	}

	updated, err := s.queries.UpdateLead(ctx, lead)
	if err != nil {
		return nil, ErrLeadNotFound
	}

	// Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicLeadUpdated, "lead_updated", map[string]interface{}{
		"id":     updated.ID,
		"email":  updated.Email,
		"status": updated.Status,
	})

	return &updated, nil
}

// DeleteLead removes a lead by ID.
func (s *LeadService) DeleteLead(ctx context.Context, id int32) error {
	if err := s.queries.DeleteLead(ctx, id); err != nil {
		return ErrLeadNotFound
	}

	// Kafka event
	_ = s.kafka.Publish(ctx, kafka.TopicLeadDeleted, "lead_deleted", map[string]interface{}{
		"id": id,
	})

	return nil
}

// GetAllLeads returns a paginated list of leads.
func (s *LeadService) GetAllLeads(ctx context.Context, pageNumber, pageSize int32) ([]db.Lead, error) {
	if pageNumber == 0 {
		pageNumber = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	offset := (pageNumber - 1) * pageSize

	return s.queries.GetAll(ctx, db.GetAllParams{
		Limit:  pageSize,
		Offset: offset,
	})
}

// GetLeadByEmail retrieves a lead by email.
func (s *LeadService) GetLeadByEmail(ctx context.Context, email string) (*db.Lead, error) {
	if strings.TrimSpace(email) == "" {
		return nil, ErrInvalidLeadData
	}
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	lead, err := s.queries.GetLeadByEmail(ctx, email)
	if err != nil {
		return nil, ErrLeadNotFound
	}
	return &lead, nil
}
