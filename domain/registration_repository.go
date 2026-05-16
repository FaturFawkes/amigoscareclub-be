package domain

import "context"

// RegistrationRepository defines storage behavior for registrations.
type RegistrationRepository interface {
	Save(ctx context.Context, reg *TicketRegistration) error
	GetByID(ctx context.Context, id RegistrationID) (*TicketRegistration, error)
}
