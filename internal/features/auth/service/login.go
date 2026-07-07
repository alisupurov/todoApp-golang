package auth_service

import (
	"context"
	"fmt"

	core_auth "github.com/alisupurov/todoApp-golang/internal/core/auth"
	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	accountID, passwordHash, err := s.authRepository.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials: %w", core_errors.ErrUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials: %w", core_errors.ErrUnauthorized)
	}

	token, err := core_auth.GenerateToken(accountID, s.jwtConfig)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
