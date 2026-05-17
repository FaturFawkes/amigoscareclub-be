package usecase

import (
	"context"

	"myapp/application/dto"
	"myapp/application/serviceInterface"
	"myapp/domain"
)

// AdminVerifyRegistrationUseCase transitions a registration's status.
type AdminVerifyRegistrationUseCase struct {
	repo          domain.RegistrationRepository
	clock         serviceInterface.Clock
	emailNotifier serviceInterface.EmailNotifier
}

// NewAdminVerifyRegistrationUseCase wires the repository, clock, and optional email notifier.
func NewAdminVerifyRegistrationUseCase(repo domain.RegistrationRepository, clock serviceInterface.Clock, emailNotifier serviceInterface.EmailNotifier) *AdminVerifyRegistrationUseCase {
	return &AdminVerifyRegistrationUseCase{repo: repo, clock: clock, emailNotifier: emailNotifier}
}

// Execute applies a status transition (verified or rejected) and persists the change.
func (uc *AdminVerifyRegistrationUseCase) Execute(ctx context.Context, input dto.VerifyRegistrationInput) (dto.VerifyRegistrationOutput, error) {
	reg, err := uc.repo.GetByID(ctx, input.EventSlug, domain.RegistrationID(input.RegistrationID))
	if err != nil {
		return dto.VerifyRegistrationOutput{}, err
	}

	now := uc.clock.Now()
	switch domain.RegistrationStatus(input.Status) {
	case domain.StatusVerified:
		if err := reg.Verify(now); err != nil {
			return dto.VerifyRegistrationOutput{}, err
		}
	case domain.StatusRejected:
		if err := reg.Reject(input.Note, now); err != nil {
			return dto.VerifyRegistrationOutput{}, err
		}
	default:
		return dto.VerifyRegistrationOutput{}, domain.ErrInvalidStatusTransition
	}

	if err := uc.repo.Update(ctx, reg); err != nil {
		return dto.VerifyRegistrationOutput{}, err
	}

	if reg.Status == domain.StatusVerified && uc.emailNotifier != nil {
		_ = uc.emailNotifier.SendVerificationConfirmation(ctx, reg.Runner.Email, reg.Runner.Name, reg.TicketNumber)
	}

	return dto.VerifyRegistrationOutput{Data: registrationToData(reg)}, nil
}
