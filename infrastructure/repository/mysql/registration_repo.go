package mysql

import (
	"context"
	"database/sql"
	"errors"

	"myapp/domain"
)

// RegistrationRepository is a MySQL implementation of the domain repository.
type RegistrationRepository struct {
	db *sql.DB
}

// NewRegistrationRepository creates a new repo instance.
func NewRegistrationRepository(db *sql.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

// Save persists the registration aggregate.
func (r *RegistrationRepository) Save(ctx context.Context, reg *domain.TicketRegistration) error {
	_ = ctx
	_ = reg
	return errors.New("not implemented")
}

// GetByID fetches a registration by its id.
func (r *RegistrationRepository) GetByID(ctx context.Context, id domain.RegistrationID) (*domain.TicketRegistration, error) {
	_ = ctx
	_ = id
	return nil, errors.New("not implemented")
}
