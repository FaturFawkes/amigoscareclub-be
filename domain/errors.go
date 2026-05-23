package domain

import "errors"

var (
	ErrEventNotFound           = errors.New("event not found")
	ErrRegistrationNotFound    = errors.New("registration not found")
	ErrAdminNotFound           = errors.New("admin not found")
	ErrDuplicateRegistration   = errors.New("duplicate registration")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidTicketStatus     = errors.New("invalid ticket status")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrFileTooLarge            = errors.New("file too large")
	ErrInvalidMIMEType         = errors.New("invalid file type")
)
