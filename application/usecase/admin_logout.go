package usecase

import (
	"context"

	"myapp/application/serviceInterface"
	"myapp/domain"
)

// AdminLogoutUseCase revokes the provided JWT token.
type AdminLogoutUseCase struct {
	tokenRepo domain.TokenRepository
	tokenSvc  serviceInterface.TokenService
}

// NewAdminLogoutUseCase wires the token repository and token service.
func NewAdminLogoutUseCase(tokenRepo domain.TokenRepository, tokenSvc serviceInterface.TokenService) *AdminLogoutUseCase {
	return &AdminLogoutUseCase{tokenRepo: tokenRepo, tokenSvc: tokenSvc}
}

// Execute parses the raw token, extracts the JTI, and blacklists it.
func (uc *AdminLogoutUseCase) Execute(ctx context.Context, rawToken string) error {
	claims, err := uc.tokenSvc.Parse(ctx, rawToken)
	if err != nil {
		return domain.ErrUnauthorized
	}
	return uc.tokenRepo.Revoke(ctx, claims.JTI, claims.ExpiresAt)
}
