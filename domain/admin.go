package domain

import "time"

// AdminID is the unique identifier for an admin account.
type AdminID string

// Admin represents an administrator who can manage event registrations.
type Admin struct {
	ID           AdminID
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
