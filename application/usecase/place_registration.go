package usecase

import (
	"context"
	"time"

	"myapp/application/dto"
	"myapp/application/serviceInterface"
	"myapp/domain"
)

// PlaceRegistrationUseCase handles registration creation.
type PlaceRegistrationUseCase struct {
	repo          domain.RegistrationRepository
	emailNotifier serviceInterface.EmailNotifier
	smsNotifier   serviceInterface.SMSNotifier
}

// NewPlaceRegistrationUseCase wires the dependencies.
func NewPlaceRegistrationUseCase(
	repo domain.RegistrationRepository,
	emailNotifier serviceInterface.EmailNotifier,
	smsNotifier serviceInterface.SMSNotifier,
) *PlaceRegistrationUseCase {
	return &PlaceRegistrationUseCase{
		repo:          repo,
		emailNotifier: emailNotifier,
		smsNotifier:   smsNotifier,
	}
}

// Execute creates a new registration aggregate and persists it.
func (uc *PlaceRegistrationUseCase) Execute(ctx context.Context, input dto.PlaceRegistrationInput) (dto.PlaceRegistrationOutput, error) {
	registration := domain.NewTicketRegistration(
		domain.RegistrationID(input.RegistrationID),
		domain.Runner{
			Name:  input.RunnerName,
			Email: input.RunnerEmail,
			Phone: input.RunnerPhone,
		},
		domain.Event{
			ID:       input.EventID,
			Name:     input.EventName,
			Date:     input.EventDate,
			Location: input.EventLocation,
		},
		domain.TicketCategory(input.Category),
		time.Now(),
	)

	if err := uc.repo.Save(ctx, registration); err != nil {
		return dto.PlaceRegistrationOutput{}, err
	}

	if uc.emailNotifier != nil {
		if err := uc.emailNotifier.SendRegistrationConfirmation(ctx, input.RunnerEmail, input.RegistrationID); err != nil {
			return dto.PlaceRegistrationOutput{}, err
		}
	}

	if uc.smsNotifier != nil {
		if err := uc.smsNotifier.SendRegistrationConfirmation(ctx, input.RunnerPhone, input.RegistrationID); err != nil {
			return dto.PlaceRegistrationOutput{}, err
		}
	}

	return dto.PlaceRegistrationOutput{RegistrationID: input.RegistrationID}, nil
}
