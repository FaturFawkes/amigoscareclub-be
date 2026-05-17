package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"myapp/application/serviceInterface"
	"myapp/infrastructure/config"
)

// JWTService implements TokenService using HS256 signed JWTs.
type JWTService struct {
	secret []byte
	ttl    time.Duration
}

// NewJWTService creates a JWTService using the application config.
func NewJWTService(cfg config.Config) *JWTService {
	return &JWTService{
		secret: []byte(cfg.JWTSecret),
		ttl:    cfg.JWTTTL,
	}
}

type jwtClaims struct {
	jwt.RegisteredClaims
}

// Issue signs a new JWT for the given admin ID.
func (s *JWTService) Issue(_ context.Context, adminID string) (string, serviceInterface.TokenClaims, error) {
	now := time.Now()
	exp := now.Add(s.ttl)
	jti := ulid.Make().String()

	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   adminID,
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", serviceInterface.TokenClaims{}, fmt.Errorf("sign token: %w", err)
	}

	return signed, serviceInterface.TokenClaims{
		Sub:       adminID,
		JTI:       jti,
		ExpiresAt: exp,
	}, nil
}

// Parse validates and extracts claims from a signed JWT string.
func (s *JWTService) Parse(_ context.Context, tokenStr string) (serviceInterface.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return serviceInterface.TokenClaims{}, err
	}

	c, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return serviceInterface.TokenClaims{}, fmt.Errorf("invalid token claims")
	}

	return serviceInterface.TokenClaims{
		Sub:       c.Subject,
		JTI:       c.ID,
		ExpiresAt: c.ExpiresAt.Time,
	}, nil
}
