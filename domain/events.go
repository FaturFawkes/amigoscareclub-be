package domain

import "time"

// DomainEvent is a marker interface for domain events.
type DomainEvent interface {
	Name() string
	OccurredAt() time.Time
}

// TicketRegistered is raised when a registration is created.
type TicketRegistered struct {
	RegistrationID RegistrationID
	Occurred       time.Time
}

func (e TicketRegistered) Name() string {
	return "ticket.registered"
}

func (e TicketRegistered) OccurredAt() time.Time {
	return e.Occurred
}
