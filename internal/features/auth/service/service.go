package auth_service

import (
	"context"

	core_auth "github.com/alisupurov/todoApp-golang/internal/core/auth"
)

type AuthService struct {
	authRepository AuthRepository
	jwtConfig      core_auth.Config
}

type AuthRepository interface {
	GetAccountByEmail(ctx context.Context, email string) (accountID int, passwordHash string, err error)
	CreateAccount(ctx context.Context, email, passwordHash string) (accountID int, err error)
}

func NewAuthService(authRepository AuthRepository, jwtConfig core_auth.Config) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		jwtConfig:      jwtConfig,
	}
}
