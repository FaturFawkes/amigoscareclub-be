package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"

	"myapp/application/dto"
	"myapp/application/serviceInterface"
	"myapp/domain"
)

// CreateRegistrationUseCase handles new event registration.
type CreateRegistrationUseCase struct {
	registrationRepo domain.RegistrationRepository
	eventRepo        domain.EventRepository
	storage          serviceInterface.FileStorage
	emailNotifier    serviceInterface.EmailNotifier
	idGen            serviceInterface.IDGenerator
	clock            serviceInterface.Clock
}

// NewCreateRegistrationUseCase wires all dependencies.
func NewCreateRegistrationUseCase(
	registrationRepo domain.RegistrationRepository,
	eventRepo domain.EventRepository,
	storage serviceInterface.FileStorage,
	emailNotifier serviceInterface.EmailNotifier,
	idGen serviceInterface.IDGenerator,
	clock serviceInterface.Clock,
) *CreateRegistrationUseCase {
	return &CreateRegistrationUseCase{
		registrationRepo: registrationRepo,
		eventRepo:        eventRepo,
		storage:          storage,
		emailNotifier:    emailNotifier,
		idGen:            idGen,
		clock:            clock,
	}
}

// Execute validates, persists the registration, and uploads the payment proof.
func (uc *CreateRegistrationUseCase) Execute(ctx context.Context, input dto.CreateRegistrationInput) (dto.CreateRegistrationOutput, error) {
	if _, err := uc.eventRepo.GetBySlug(ctx, input.EventSlug); err != nil {
		return dto.CreateRegistrationOutput{}, err
	}

	existing, err := uc.registrationRepo.FindByEventAndEmail(ctx, input.EventSlug, input.Email)
	if err != nil && !errors.Is(err, domain.ErrRegistrationNotFound) {
		return dto.CreateRegistrationOutput{}, err
	}
	if existing != nil {
		return dto.CreateRegistrationOutput{}, domain.ErrDuplicateRegistration
	}

	id, err := uc.idGen.NewRegistrationID(ctx)
	if err != nil {
		return dto.CreateRegistrationOutput{}, err
	}
	ticketNumber, err := uc.idGen.NewTicketNumber(ctx)
	if err != nil {
		return dto.CreateRegistrationOutput{}, err
	}

	now := uc.clock.Now()
	runner := domain.Runner{
		Name:         input.Name,
		Email:        input.Email,
		Phone:        input.Phone,
		Age:          input.Age,
		CoffeeChoice: input.CoffeeChoice,
	}
	reg := domain.NewRegistration(domain.RegistrationID(id), ticketNumber, input.EventSlug, runner, now)

	if input.PaymentProof != nil {
		f, err := input.PaymentProof.Open()
		if err != nil {
			return dto.CreateRegistrationOutput{}, fmt.Errorf("open payment proof: %w", err)
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return dto.CreateRegistrationOutput{}, fmt.Errorf("read payment proof: %w", err)
		}
		contentType := input.PaymentProof.Header.Get("Content-Type")
		key := fmt.Sprintf("proofs/%s", id)
		if err := uc.storage.Put(ctx, key, data, contentType); err != nil {
			return dto.CreateRegistrationOutput{}, err
		}
		url, err := uc.storage.GetURL(ctx, key)
		if err != nil {
			return dto.CreateRegistrationOutput{}, err
		}
		reg.PaymentProofURL = url
	}

	if err := uc.registrationRepo.Save(ctx, reg); err != nil {
		return dto.CreateRegistrationOutput{}, err
	}

	if uc.emailNotifier != nil {
		_ = uc.emailNotifier.SendRegistrationConfirmation(ctx, input.Email, id)
	}

	out := dto.CreateRegistrationOutput{Data: registrationToData(reg)}
	out.Meta.Message = "Pendaftaran berhasil! Tim kami akan memverifikasi pembayaranmu dan mengirimkan email konfirmasi tiket setelah diverifikasi oleh admin."
	return out, nil
}
