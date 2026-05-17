package domain

import "context"

// RegistrationFilter holds optional filter criteria for listing registrations.
type RegistrationFilter struct {
	Status *RegistrationStatus
}

// RegistrationRepository defines persistence operations for registrations.
type RegistrationRepository interface {
	Save(ctx context.Context, reg *Registration) error
	GetByID(ctx context.Context, eventSlug string, id RegistrationID) (*Registration, error)
	FindByEventAndEmail(ctx context.Context, eventSlug, email string) (*Registration, error)
	List(ctx context.Context, eventSlug string, filter RegistrationFilter, page, perPage int) ([]*Registration, int, error)
	Update(ctx context.Context, reg *Registration) error
}
