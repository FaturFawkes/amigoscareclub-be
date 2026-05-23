package usecase

import (
	"context"
	"fmt"

	"myapp/application/dto"
	"myapp/application/serviceInterface"
	"myapp/domain"
)

// AdminResendTicketUseCase resends the ticket email for a single registration.
type AdminResendTicketUseCase struct {
	regRepo       domain.RegistrationRepository
	eventRepo     domain.EventRepository
	clock         serviceInterface.Clock
	emailNotifier serviceInterface.EmailNotifier
}

// NewAdminResendTicketUseCase wires dependencies.
func NewAdminResendTicketUseCase(
	regRepo domain.RegistrationRepository,
	eventRepo domain.EventRepository,
	clock serviceInterface.Clock,
	emailNotifier serviceInterface.EmailNotifier,
) *AdminResendTicketUseCase {
	return &AdminResendTicketUseCase{
		regRepo:       regRepo,
		eventRepo:     eventRepo,
		clock:         clock,
		emailNotifier: emailNotifier,
	}
}

// Execute resends the ticket email and updates ticket_sent_at.
func (uc *AdminResendTicketUseCase) Execute(ctx context.Context, input dto.ResendTicketInput) (dto.ResendTicketOutput, error) {
	reg, err := uc.regRepo.GetByID(ctx, input.EventSlug, domain.RegistrationID(input.RegistrationID))
	if err != nil {
		return dto.ResendTicketOutput{}, err
	}

	now := uc.clock.Now()
	if err := reg.MarkTicketResent(now); err != nil {
		return dto.ResendTicketOutput{}, err
	}

	event, err := uc.eventRepo.GetBySlug(ctx, input.EventSlug)
	if err != nil {
		return dto.ResendTicketOutput{}, err
	}

	if err := uc.regRepo.Update(ctx, reg); err != nil {
		return dto.ResendTicketOutput{}, err
	}

	if uc.emailNotifier != nil {
		eventDate := fmt.Sprintf("%s %d %s %d",
			weekdayID(event.Date.Weekday()),
			event.Date.Day(),
			monthID(event.Date.Month()),
			event.Date.Year(),
		)
		_ = uc.emailNotifier.SendTicket(ctx,
			reg.Runner.Email, reg.Runner.Name, reg.TicketNumber,
			event.Title, eventDate, event.Time, event.Location,
		)
	}

	var out dto.ResendTicketOutput
	out.Data.Email = reg.Runner.Email
	return out, nil
}
