package dto

import (
	"mime/multipart"
	"time"
)

// CreateRegistrationInput holds data from a multipart registration form.
type CreateRegistrationInput struct {
	EventSlug    string
	Name         string
	Email        string
	Phone        string
	Age          int
	CoffeeChoice string
	PaymentProof *multipart.FileHeader
}

// RegistrationData matches the swagger Registration schema.
type RegistrationData struct {
	ID              string     `json:"id"`
	TicketNumber    string     `json:"ticket_number"`
	EventSlug       string     `json:"event_slug"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Age             int        `json:"age"`
	CoffeeChoice    string     `json:"coffee_choice"`
	Status          string     `json:"status"`
	Note            *string    `json:"note"`
	PaymentProofURL *string    `json:"payment_proof_url"`
	RegisteredAt    time.Time  `json:"registered_at"`
	VerifiedAt      *time.Time `json:"verified_at"`
	TicketSentAt    *time.Time `json:"ticket_sent_at"`
}

// CreateRegistrationOutput matches the swagger 201 response for registration.
type CreateRegistrationOutput struct {
	Data RegistrationData `json:"data"`
	Meta struct {
		Message string `json:"message"`
	} `json:"meta"`
}

// GetRegistrationInput carries identifiers to look up a single registration.
type GetRegistrationInput struct {
	EventSlug      string
	RegistrationID string
}

// GetRegistrationOutput wraps a single registration.
type GetRegistrationOutput struct {
	Data RegistrationData `json:"data"`
}

// ListRegistrationsInput carries filter and pagination parameters.
type ListRegistrationsInput struct {
	EventSlug string
	Status    string
	Page      int
	PerPage   int
}

// PaginationMeta matches the swagger PaginationMeta schema.
type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
}

// ListRegistrationsOutput wraps a paginated list of registrations.
type ListRegistrationsOutput struct {
	Data []RegistrationData `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

// VerifyRegistrationInput carries the status transition request from admin.
type VerifyRegistrationInput struct {
	EventSlug      string
	RegistrationID string
	Status         string
	Note           string
}

// VerifyRegistrationOutput wraps the updated registration after verification.
type VerifyRegistrationOutput struct {
	Data RegistrationData `json:"data"`
}

// ResendTicketInput carries identifiers for the resend-ticket operation.
type ResendTicketInput struct {
	EventSlug      string
	RegistrationID string
}

// ResendTicketOutput holds the recipient email address.
type ResendTicketOutput struct {
	Data struct {
		Email string `json:"email"`
	} `json:"data"`
}

// ResendAllTicketsInput carries the event slug for bulk resend.
type ResendAllTicketsInput struct {
	EventSlug string
}

// ResendAllTicketsOutput holds the number of emails dispatched.
type ResendAllTicketsOutput struct {
	Data struct {
		Sent int `json:"sent"`
	} `json:"data"`
}
