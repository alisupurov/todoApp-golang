package auth_service

import (
	"context"
	"fmt"

	core_auth "github.com/alisupurov/todoApp-golang/internal/core/auth"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Register(ctx context.Context, email, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	accountID, err := s.authRepository.CreateAccount(ctx, email, string(hash))
	if err != nil {
		return "", fmt.Errorf("create account: %w", err)
	}

	token, err := core_auth.GenerateToken(accountID, s.jwtConfig)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
