package usecase

import (
	"context"

	"myapp/application/dto"
	"myapp/domain"
)

// AdminMeUseCase retrieves the profile of the currently authenticated admin.
type AdminMeUseCase struct {
	adminRepo domain.AdminRepository
}

// NewAdminMeUseCase wires the repository.
func NewAdminMeUseCase(adminRepo domain.AdminRepository) *AdminMeUseCase {
	return &AdminMeUseCase{adminRepo: adminRepo}
}

// Execute returns the admin profile for the given admin ID.
func (uc *AdminMeUseCase) Execute(ctx context.Context, adminID string) (dto.AdminMeOutput, error) {
	admin, err := uc.adminRepo.GetByID(ctx, domain.AdminID(adminID))
	if err != nil {
		return dto.AdminMeOutput{}, err
	}
	return dto.AdminMeOutput{
		Data: dto.AdminProfileData{
			ID:    string(admin.ID),
			Name:  admin.Name,
			Email: admin.Email,
		},
	}, nil
}
