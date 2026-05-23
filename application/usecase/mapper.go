package usecase

import (
	"time"

	"myapp/application/dto"
	"myapp/domain"
)

func weekdayID(w time.Weekday) string {
	names := [...]string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	return names[int(w)%7]
}

func monthID(m time.Month) string {
	names := [...]string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	if int(m) < 1 || int(m) > 12 {
		return ""
	}
	return names[int(m)]
}

func registrationToData(reg *domain.Registration) dto.RegistrationData {
	var proofURL *string
	if reg.PaymentProofURL != "" {
		s := reg.PaymentProofURL
		proofURL = &s
	}
	var note *string
	if reg.Note != "" {
		s := reg.Note
		note = &s
	}
	return dto.RegistrationData{
		ID:              string(reg.ID),
		TicketNumber:    reg.TicketNumber,
		EventSlug:       reg.EventSlug,
		Name:            reg.Runner.Name,
		Email:           reg.Runner.Email,
		Phone:           reg.Runner.Phone,
		Age:             reg.Runner.Age,
		CoffeeChoice:    reg.Runner.CoffeeChoice,
		Status:          string(reg.Status),
		Note:            note,
		PaymentProofURL: proofURL,
		RegisteredAt:    reg.RegisteredAt,
		VerifiedAt:      reg.VerifiedAt,
		TicketSentAt:    reg.TicketSentAt,
	}
}
