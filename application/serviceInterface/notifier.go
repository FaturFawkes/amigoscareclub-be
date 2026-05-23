package serviceInterface

import "context"

// EmailNotifier sends registration notifications by email.
type EmailNotifier interface {
	SendRegistrationConfirmation(ctx context.Context, email, name, registrationID string) error
	SendVerificationConfirmation(ctx context.Context, email, name, ticketNumber string) error
	SendRejectionNotification(ctx context.Context, email, name, note string) error
	SendTicket(ctx context.Context, email, name, ticketNumber, eventTitle, eventDate, eventTime, eventLocation string) error
}

// SMSNotifier sends registration notifications by SMS.
type SMSNotifier interface {
	SendRegistrationConfirmation(ctx context.Context, phone string, registrationID string) error
}
