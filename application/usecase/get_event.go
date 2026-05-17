package usecase

import (
	"context"

	"myapp/application/dto"
	"myapp/domain"
)

// GetEventUseCase retrieves an event by slug.
type GetEventUseCase struct {
	repo domain.EventRepository
}

// NewGetEventUseCase wires the repository.
func NewGetEventUseCase(repo domain.EventRepository) *GetEventUseCase {
	return &GetEventUseCase{repo: repo}
}

// Execute returns the event details as a DTO.
func (uc *GetEventUseCase) Execute(ctx context.Context, input dto.GetEventInput) (dto.GetEventOutput, error) {
	event, err := uc.repo.GetBySlug(ctx, input.Slug)
	if err != nil {
		return dto.GetEventOutput{}, err
	}
	return dto.GetEventOutput{
		Data: dto.EventData{
			Slug:             event.Slug,
			Title:            event.Title,
			Date:             event.Date.Format("2006-01-02"),
			Time:             event.Time,
			Timezone:         event.Timezone,
			Location:         event.Location,
			DistanceKm:       event.DistanceKm,
			Pace:             event.Pace,
			RegistrationOpen: event.RegistrationOpen,
			CoffeeOptions:    event.CoffeeOptions,
			Payment: dto.EventPaymentData{
				Bank:          event.Payment.Bank,
				AccountNumber: event.Payment.AccountNumber,
				AccountName:   event.Payment.AccountName,
			},
		},
	}, nil
}
