package services

import (
	"context"
	"crm/internal/adapters/database/db"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/segmentio/kafka-go"
)

var (
	ErrContactNotFound    = errors.New("contact not found")
	ErrInvalidContactData = errors.New("invalid contact data")
	ErrContactExists      = errors.New("contact with this email already exists")
)

type ContactServiceInterface interface {
	CreateContact(ctx context.Context, contact db.CreateContactParams) (*db.Contact, error)
	GetContact(ctx context.Context, id int32) (*db.Contact, error)
	UpdateContact(ctx context.Context, contact db.UpdateContactParams) (*db.Contact, error)
	DeleteContact(ctx context.Context, id int32) error
	ListContacts(ctx context.Context, pageNumber, pageSize int32) ([]db.Contact, error)
}

type ContactService struct {
	queries *db.Queries
	kafka   *kafka.Writer
}

func NewContactService(queries *db.Queries, kafkaWriter *kafka.Writer) *ContactService {
	return &ContactService{queries: queries, kafka: kafkaWriter}
}

// CreateContact validates and creates a new unified contact.
// It ensures that required fields are present based on the ContactType.
func (s *ContactService) CreateContact(ctx context.Context, contact db.CreateContactParams) (*db.Contact, error) {
	// Email (string, NOT NULL)
	if strings.TrimSpace(contact.Email) == "" {
		return nil, ErrInvalidContactData
	}
	if !isValidEmail(contact.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate based on type
	switch contact.ContactType {
	case "individual":
		// sql.NullString must be checked with .Valid
		if !contact.FirstName.Valid || !contact.LastName.Valid {
			return nil, ErrInvalidContactData
		}
	case "company":
		if !contact.CompanyName.Valid {
			return nil, ErrInvalidContactData
		}
	default:
		return nil, errors.New("unknown contact type")
	}

	// Insert into DB
	createdContact, err := s.queries.CreateContact(ctx, contact)
	if err != nil {
		return nil, err
	}

	// Kafka event
	if err := s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte("contact_created"),
		Value: []byte(createdContact.Email), // Email is plain string
	}); err != nil {
		log.Printf("failed to write kafka message: %v", err)
	}

	return &createdContact, nil
}

// GetContact retrieves a contact by its ID.
func (s *ContactService) GetContact(ctx context.Context, id int32) (*db.Contact, error) {
	contact, err := s.queries.GetContact(ctx, id)
	if err != nil {
		return nil, ErrContactNotFound
	}
	return &contact, nil
}

// UpdateContact validates and updates an existing contact.
func (s *ContactService) UpdateContact(ctx context.Context, contact db.UpdateContactParams) (*db.Contact, error) {
	if contact.ID == 0 {
		return nil, ErrInvalidContactData
	}
	if contact.Email != "" && !isValidEmail(contact.Email) {
		return nil, errors.New("invalid email format")
	}

	updatedContact, err := s.queries.UpdateContact(ctx, contact)
	if err != nil {
		return nil, ErrContactNotFound
	}

	// Kafka event
	err = s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte("contact_updated"),
		Value: []byte(updatedContact.Email),
	})
	if err != nil {
		log.Printf("failed to write kafka message: %v", err)
	}

	return &updatedContact, nil
}

// DeleteContact removes a contact by its ID.
func (s *ContactService) DeleteContact(ctx context.Context, id int32) error {
	err := s.queries.DeleteContact(ctx, id)
	if err != nil {
		return ErrContactNotFound
	}

	// Kafka event
	err = s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte("contact_deleted"),
		Value: []byte(string(rune(id))),
	})
	if err != nil {
		log.Printf("failed to write kafka message: %v", err)
	}

	return nil
}

// ListContacts retrieves contacts with pagination and sorting.
func (s *ContactService) ListContacts(ctx context.Context, pageNumber, pageSize int32) ([]db.Contact, error) {
	if pageNumber == 0 {
		pageNumber = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	offset := (pageNumber - 1) * pageSize

	contacts, err := s.queries.ListContacts(ctx, db.ListContactsParams{
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

// isValidEmail validates the email format using a regular expression.
func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
