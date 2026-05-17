package usecase

import (
	"context"
	"errors"

	"myapp/application/dto"
	"myapp/application/serviceInterface"
	"myapp/domain"
)

// AdminLoginUseCase authenticates an admin and issues a JWT.
type AdminLoginUseCase struct {
	adminRepo domain.AdminRepository
	hasher    serviceInterface.PasswordHasher
	tokenSvc  serviceInterface.TokenService
}

// NewAdminLoginUseCase wires all dependencies.
func NewAdminLoginUseCase(
	adminRepo domain.AdminRepository,
	hasher serviceInterface.PasswordHasher,
	tokenSvc serviceInterface.TokenService,
) *AdminLoginUseCase {
	return &AdminLoginUseCase{adminRepo: adminRepo, hasher: hasher, tokenSvc: tokenSvc}
}

// Execute validates email + password and returns a Bearer token on success.
func (uc *AdminLoginUseCase) Execute(ctx context.Context, input dto.LoginInput) (dto.LoginOutput, error) {
	admin, err := uc.adminRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			return dto.LoginOutput{}, domain.ErrInvalidCredentials
		}
		return dto.LoginOutput{}, err
	}

	if err := uc.hasher.Compare(ctx, admin.PasswordHash, input.Password); err != nil {
		return dto.LoginOutput{}, domain.ErrInvalidCredentials
	}

	token, claims, err := uc.tokenSvc.Issue(ctx, string(admin.ID))
	if err != nil {
		return dto.LoginOutput{}, err
	}

	return dto.LoginOutput{
		Data: dto.LoginData{
			Token:     token,
			ExpiresAt: claims.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
			Admin: dto.AdminProfileData{
				ID:    string(admin.ID),
				Name:  admin.Name,
				Email: admin.Email,
			},
		},
	}, nil
}
