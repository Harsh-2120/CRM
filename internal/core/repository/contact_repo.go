package repository

import (
	"context"
	"crm/internal/adapters/database/db"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
)

var (
	ErrContactExists   = errors.New("contact with this email already exists")
	ErrContactNotFound = errors.New("contact not found")
)

// ContactRepository defines the contract for contact CRUD operations.
type ContactRepository interface {
	Create(ctx context.Context, arg db.CreateContactParams) (db.Contact, error)
	GetByID(ctx context.Context, id int32) (db.Contact, error)
	Update(ctx context.Context, arg db.UpdateContactParams) (db.Contact, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, limit, offset int32, sortBy string, ascending bool) ([]db.Contact, error)
}

type contactRepository struct {
	q *db.Queries
}

// NewContactRepository creates a new instance of contactRepository.
func NewContactRepository(q *db.Queries) ContactRepository {
	return &contactRepository{q: q}
}

// ----------------- CRUD -----------------

func (r *contactRepository) Create(ctx context.Context, arg db.CreateContactParams) (db.Contact, error) {
	contact, err := r.q.CreateContact(ctx, arg)
	if err != nil {
		if isUniqueConstraintError(err, "contacts_email_key") {
			return db.Contact{}, ErrContactExists
		}
		return db.Contact{}, err
	}
	return contact, nil
}

func (r *contactRepository) GetByID(ctx context.Context, id int32) (db.Contact, error) {
	contact, err := r.q.GetContact(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Contact{}, ErrContactNotFound
		}
		return db.Contact{}, err
	}
	return contact, nil
}

func (r *contactRepository) Update(ctx context.Context, arg db.UpdateContactParams) (db.Contact, error) {
	contact, err := r.q.UpdateContact(ctx, arg)
	if err != nil {
		if isUniqueConstraintError(err, "contacts_email_key") {
			return db.Contact{}, ErrContactExists
		}
		if errors.Is(err, sql.ErrNoRows) {
			return db.Contact{}, ErrContactNotFound
		}
		return db.Contact{}, err
	}
	return contact, nil
}

func (r *contactRepository) Delete(ctx context.Context, id int32) error {
	err := r.q.DeleteContact(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrContactNotFound
		}
		return err
	}
	return nil
}

func (r *contactRepository) List(ctx context.Context, limit, offset int32, sortBy string, ascending bool) ([]db.Contact, error) {
	// sqlc normally generates `ListContacts(limit, offset)` without sorting.
	// If you need sorting, modify your SQL in `contacts.sql` to accept ORDER BY dynamically.
	contacts, err := r.q.ListContacts(ctx, db.ListContactsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

// ----------------- Helpers -----------------

// isUniqueConstraintError checks for PostgreSQL unique constraint violations.
func isUniqueConstraintError(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" && pgErr.ConstraintName == constraintName {
			return true
		}
	}
	return false
}
