package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"

	"myapp/application/dto"
	"myapp/application/serviceInterface"
	"myapp/domain"
)

// AdminResendAllTicketsUseCase resends ticket emails to all eligible registrations for an event.
type AdminResendAllTicketsUseCase struct {
	regRepo       domain.RegistrationRepository
	eventRepo     domain.EventRepository
	clock         serviceInterface.Clock
	emailNotifier serviceInterface.EmailNotifier
}

// NewAdminResendAllTicketsUseCase wires dependencies.
func NewAdminResendAllTicketsUseCase(
	regRepo domain.RegistrationRepository,
	eventRepo domain.EventRepository,
	clock serviceInterface.Clock,
	emailNotifier serviceInterface.EmailNotifier,
) *AdminResendAllTicketsUseCase {
	return &AdminResendAllTicketsUseCase{
		regRepo:       regRepo,
		eventRepo:     eventRepo,
		clock:         clock,
		emailNotifier: emailNotifier,
	}
}

// Execute dispatches ticket emails to all verified/ticket_sent registrations asynchronously.
func (uc *AdminResendAllTicketsUseCase) Execute(ctx context.Context, input dto.ResendAllTicketsInput) (dto.ResendAllTicketsOutput, error) {
	event, err := uc.eventRepo.GetBySlug(ctx, input.EventSlug)
	if err != nil {
		return dto.ResendAllTicketsOutput{}, err
	}

	regs, err := uc.regRepo.ListEligibleForTicket(ctx, input.EventSlug)
	if err != nil {
		return dto.ResendAllTicketsOutput{}, err
	}

	eventDate := fmt.Sprintf("%s %d %s %d",
		weekdayID(event.Date.Weekday()),
		event.Date.Day(),
		monthID(event.Date.Month()),
		event.Date.Year(),
	)

	var wg sync.WaitGroup
	for _, reg := range regs {
		wg.Add(1)
		go func(r *domain.Registration) {
			defer wg.Done()
			bgCtx := context.Background()

			now := uc.clock.Now()
			if err := r.MarkTicketResent(now); err != nil {
				log.Printf("resend_all: skip registration %s: %v", r.ID, err)
				return
			}

			if err := uc.regRepo.Update(bgCtx, r); err != nil {
				log.Printf("resend_all: update registration %s: %v", r.ID, err)
				return
			}

			if uc.emailNotifier != nil {
				if err := uc.emailNotifier.SendTicket(bgCtx,
					r.Runner.Email, r.Runner.Name, r.TicketNumber,
					event.Title, eventDate, event.Time, event.Location,
				); err != nil {
					log.Printf("resend_all: send ticket to %s: %v", r.Runner.Email, err)
				}
			}
		}(reg)
	}

	var out dto.ResendAllTicketsOutput
	out.Data.Sent = len(regs)

	go func() { wg.Wait() }()

	return out, nil
}
