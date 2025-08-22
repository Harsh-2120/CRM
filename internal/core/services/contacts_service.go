package services

import (
	"crm/internal/core/domain/models"
	"crm/internal/core/repository"
	"errors"
	"time"
)

var (
	ErrContactNotFound    = errors.New("contact not found")
	ErrInvalidContactData = errors.New("invalid contact data")
	ErrContactExists      = errors.New("contact with this email already exists")
)

type ContactService interface {
	CreateContact(contact *models.Contact) (*models.Contact, error)
	GetContact(id uint) (*models.Contact, error)
	UpdateContact(contact *models.Contact) (*models.Contact, error)
	DeleteContact(id uint) error
	ListContacts(pageNumber uint, pageSize uint, sortBy string, ascending bool) ([]models.Contact, error)
}

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{repo: repo}
}

// CreateContact validates and creates a new unified contact.
// It ensures that required fields are present based on the ContactType.
func (s *contactService) CreateContact(contact *models.Contact) (*models.Contact, error) {
	// Email is always required.
	if contact.Email == "" {
		return nil, ErrInvalidContactData
	}

	// Validate email format.
	if !isValidEmail(contact.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate required fields based on contact type.
	switch contact.ContactType {
	case "individual":
		if contact.FirstName == "" || contact.LastName == "" {
			return nil, ErrInvalidContactData
		}
	case "company":
		if contact.CompanyName == "" {
			return nil, ErrInvalidContactData
		}
	default:
		return nil, errors.New("unknown contact type")
	}

	// Set creation and update timestamps.
	now := time.Now()
	contact.CreatedAt = now
	contact.UpdatedAt = now

	// Attempt to create the contact in the repository.
	createdContact, err := s.repo.Create(contact)
	if err != nil {
		if errors.Is(err, repository.ErrContactExists) {
			return nil, ErrContactExists
		}
		return nil, err
	}

	return createdContact, nil
}

// GetContact retrieves a contact by its ID.
func (s *contactService) GetContact(id uint) (*models.Contact, error) {
	contact, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrContactNotFound) {
			return nil, ErrContactNotFound
		}
		return nil, err
	}
	return contact, nil
}

// UpdateContact validates and updates an existing contact.
func (s *contactService) UpdateContact(contact *models.Contact) (*models.Contact, error) {
	// Validate that the contact has a valid ID.
	if contact.ID == 0 {
		return nil, ErrInvalidContactData
	}

	// Validate email format if an email is provided.
	if contact.Email != "" && !isValidEmail(contact.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate required fields based on contact type.
	switch contact.ContactType {
	case "individual":
		if contact.FirstName == "" || contact.LastName == "" {
			return nil, ErrInvalidContactData
		}
	case "company":
		if contact.CompanyName == "" {
			return nil, ErrInvalidContactData
		}
	default:
		return nil, errors.New("unknown contact type")
	}

	// Update the updated_at timestamp.
	contact.UpdatedAt = time.Now()

	// Attempt to update the contact in the repository.
	updatedContact, err := s.repo.Update(contact)
	if err != nil {
		if errors.Is(err, repository.ErrContactNotFound) {
			return nil, ErrContactNotFound
		}
		if errors.Is(err, repository.ErrContactExists) {
			return nil, ErrContactExists
		}
		return nil, err
	}

	return updatedContact, nil
}

// DeleteContact removes a contact by its ID.
func (s *contactService) DeleteContact(id uint) error {
	// Verify if the contact exists.
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrContactNotFound) {
			return ErrContactNotFound
		}
		return err
	}

	// Delete the contact.
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

// ListContacts retrieves contacts with pagination and sorting.
// Extended sort fields now include "company_name" and "contact_type".
func (s *contactService) ListContacts(pageNumber uint, pageSize uint, sortBy string, ascending bool) ([]models.Contact, error) {
	// Default pagination if invalid values provided.
	if pageNumber == 0 {
		pageNumber = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	// Valid sort fields.
	validSortFields := map[string]bool{
		"first_name":   true,
		"last_name":    true,
		"email":        true,
		"company_name": true,
		"contact_type": true,
		"created_at":   true,
		"updated_at":   true,
	}

	if sortBy != "" && !validSortFields[sortBy] {
		return nil, errors.New("invalid sort field")
	}

	contacts, err := s.repo.List(pageNumber, pageSize, sortBy, ascending)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

// isValidEmail validates the email format using a regular expression.
// func isValidEmail(email string) bool {
// 	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
// 	re := regexp.MustCompile(regex)
// 	return re.MatchString(email)
// }
