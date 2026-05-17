package usecase

import (
	"context"

	"myapp/application/dto"
	"myapp/domain"
)

// AdminListRegistrationsUseCase returns a filtered, paginated list of registrations.
type AdminListRegistrationsUseCase struct {
	repo domain.RegistrationRepository
}

// NewAdminListRegistrationsUseCase wires the repository.
func NewAdminListRegistrationsUseCase(repo domain.RegistrationRepository) *AdminListRegistrationsUseCase {
	return &AdminListRegistrationsUseCase{repo: repo}
}

// Execute applies optional status filter and pagination, returning data + meta.
func (uc *AdminListRegistrationsUseCase) Execute(ctx context.Context, input dto.ListRegistrationsInput) (dto.ListRegistrationsOutput, error) {
	var filter domain.RegistrationFilter
	if input.Status != "" {
		s := domain.RegistrationStatus(input.Status)
		filter.Status = &s
	}

	page := input.Page
	if page < 1 {
		page = 1
	}
	perPage := input.PerPage
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	regs, total, err := uc.repo.List(ctx, input.EventSlug, filter, page, perPage)
	if err != nil {
		return dto.ListRegistrationsOutput{}, err
	}

	data := make([]dto.RegistrationData, len(regs))
	for i, r := range regs {
		data[i] = registrationToData(r)
	}

	totalPages := (total + perPage - 1) / perPage

	return dto.ListRegistrationsOutput{
		Data: data,
		Meta: dto.PaginationMeta{
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		},
	}, nil
}
