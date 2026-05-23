package domain

import "time"

// RegistrationID is the unique identifier for a registration.
type RegistrationID string

// RegistrationStatus represents the lifecycle state of a registration.
type RegistrationStatus string

const (
	StatusPendingVerification RegistrationStatus = "pending_verification"
	StatusVerified            RegistrationStatus = "verified"
	StatusRejected            RegistrationStatus = "rejected"
	StatusTicketSent          RegistrationStatus = "ticket_sent"
)

// Runner holds the personal data of the event participant.
type Runner struct {
	Name         string
	Email        string
	Phone        string
	Age          int
	CoffeeChoice string
}

// Registration is the aggregate root for an event registration.
type Registration struct {
	ID              RegistrationID
	TicketNumber    string
	EventSlug       string
	Runner          Runner
	Status          RegistrationStatus
	PaymentProofURL string
	Note            string
	RegisteredAt    time.Time
	VerifiedAt      *time.Time
	TicketSentAt    *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewRegistration creates a new pending registration.
func NewRegistration(id RegistrationID, ticketNumber, eventSlug string, runner Runner, now time.Time) *Registration {
	return &Registration{
		ID:           id,
		TicketNumber: ticketNumber,
		EventSlug:    eventSlug,
		Runner:       runner,
		Status:       StatusPendingVerification,
		RegisteredAt: now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Verify transitions the registration to verified status.
func (r *Registration) Verify(now time.Time) error {
	if r.Status != StatusPendingVerification {
		return ErrInvalidStatusTransition
	}
	r.Status = StatusVerified
	r.VerifiedAt = &now
	r.UpdatedAt = now
	return nil
}

// Reject transitions the registration to rejected status with an optional note.
func (r *Registration) Reject(note string, now time.Time) error {
	if r.Status != StatusPendingVerification {
		return ErrInvalidStatusTransition
	}
	r.Status = StatusRejected
	r.Note = note
	r.UpdatedAt = now
	return nil
}

// MarkTicketSent transitions a verified registration to ticket_sent.
func (r *Registration) MarkTicketSent(now time.Time) error {
	if r.Status != StatusVerified {
		return ErrInvalidStatusTransition
	}
	r.Status = StatusTicketSent
	r.TicketSentAt = &now
	r.UpdatedAt = now
	return nil
}

// MarkTicketResent updates ticket_sent_at and sets status to ticket_sent.
// Valid for registrations in verified or ticket_sent status.
func (r *Registration) MarkTicketResent(now time.Time) error {
	if r.Status != StatusVerified && r.Status != StatusTicketSent {
		return ErrInvalidTicketStatus
	}
	r.Status = StatusTicketSent
	r.TicketSentAt = &now
	r.UpdatedAt = now
	return nil
}
