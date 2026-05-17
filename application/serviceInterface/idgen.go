package serviceInterface

import "context"

// IDGenerator generates unique identifiers for domain entities.
type IDGenerator interface {
	NewRegistrationID(ctx context.Context) (string, error)
	NewAdminID(ctx context.Context) (string, error)
	NewTicketNumber(ctx context.Context) (string, error)
}
