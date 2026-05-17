package serviceInterface

import "context"

// EmailNotifier sends registration notifications by email.
type EmailNotifier interface {
	SendRegistrationConfirmation(ctx context.Context, email string, registrationID string) error
	SendVerificationConfirmation(ctx context.Context, email, name, ticketNumber string) error
}

// SMSNotifier sends registration notifications by SMS.
type SMSNotifier interface {
	SendRegistrationConfirmation(ctx context.Context, phone string, registrationID string) error
}
