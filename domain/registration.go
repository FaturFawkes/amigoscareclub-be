package domain

import "time"

// RegistrationID is the identifier for a ticket registration aggregate.
type RegistrationID string

// Runner is the attendee registering for the event.
type Runner struct {
	Name  string
	Email string
	Phone string
}

// Event is the running event being registered for.
type Event struct {
	ID       string
	Name     string
	Date     time.Time
	Location string
}

// TicketCategory is a value object for the race category.
type TicketCategory string

const (
	Category5K           TicketCategory = "5K"
	Category10K          TicketCategory = "10K"
	CategoryHalfMarathon TicketCategory = "HALF_MARATHON"
)

// TicketRegistration is the aggregate root.
type TicketRegistration struct {
	ID        RegistrationID
	Runner    Runner
	Event     Event
	Category  TicketCategory
	Paid      bool
	CreatedAt time.Time
}

// NewTicketRegistration creates a new registration aggregate.
func NewTicketRegistration(id RegistrationID, runner Runner, event Event, category TicketCategory, createdAt time.Time) *TicketRegistration {
	return &TicketRegistration{
		ID:        id,
		Runner:    runner,
		Event:     event,
		Category:  category,
		CreatedAt: createdAt,
	}
}

// MarkPaid updates the payment status.
func (r *TicketRegistration) MarkPaid() {
	r.Paid = true
}
