package usecase

import (
	"context"

	"myapp/application/dto"
	"myapp/domain"
)

// GetRegistrationUseCase retrieves a single registration by event + ID.
type GetRegistrationUseCase struct {
	repo domain.RegistrationRepository
}

// NewGetRegistrationUseCase wires the repository.
func NewGetRegistrationUseCase(repo domain.RegistrationRepository) *GetRegistrationUseCase {
	return &GetRegistrationUseCase{repo: repo}
}

// Execute fetches the registration and returns it as a DTO.
func (uc *GetRegistrationUseCase) Execute(ctx context.Context, input dto.GetRegistrationInput) (dto.GetRegistrationOutput, error) {
	reg, err := uc.repo.GetByID(ctx, input.EventSlug, domain.RegistrationID(input.RegistrationID))
	if err != nil {
		return dto.GetRegistrationOutput{}, err
	}
	return dto.GetRegistrationOutput{Data: registrationToData(reg)}, nil
}
