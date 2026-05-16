package dto

import "time"

// PlaceRegistrationInput captures request data from the delivery layer.
type PlaceRegistrationInput struct {
	RegistrationID string
	RunnerName     string
	RunnerEmail    string
	RunnerPhone    string
	EventID        string
	EventName      string
	EventDate      time.Time
	EventLocation  string
	Category       string
}
