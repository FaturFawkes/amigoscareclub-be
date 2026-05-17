package domain

import "context"

// AdminRepository defines persistence operations for admin accounts.
type AdminRepository interface {
	GetByEmail(ctx context.Context, email string) (*Admin, error)
	GetByID(ctx context.Context, id AdminID) (*Admin, error)
}
