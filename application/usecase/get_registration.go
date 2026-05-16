package usecase

import (
	"context"

	"myapp/application/dto"
	"myapp/domain"
)

// GetRegistrationUseCase fetches a registration by id.
type GetRegistrationUseCase struct {
	repo domain.RegistrationRepository
}

// NewGetRegistrationUseCase wires dependencies.
func NewGetRegistrationUseCase(repo domain.RegistrationRepository) *GetRegistrationUseCase {
	return &GetRegistrationUseCase{repo: repo}
}

// Execute returns the registration details.
func (uc *GetRegistrationUseCase) Execute(ctx context.Context, id string) (dto.PlaceRegistrationInput, error) {
	reg, err := uc.repo.GetByID(ctx, domain.RegistrationID(id))
	if err != nil {
		return dto.PlaceRegistrationInput{}, err
	}

	return dto.PlaceRegistrationInput{
		RegistrationID: string(reg.ID),
		RunnerName:     reg.Runner.Name,
		RunnerEmail:    reg.Runner.Email,
		RunnerPhone:    reg.Runner.Phone,
		EventID:        reg.Event.ID,
		EventName:      reg.Event.Name,
		EventDate:      reg.Event.Date,
		EventLocation:  reg.Event.Location,
		Category:       string(reg.Category),
	}, nil
}
