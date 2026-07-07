package core_auth

import (
	"fmt"
	"time"

	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AccountID int `json:"account_id"`
	jwt.RegisteredClaims
}

func GenerateToken(accountID int, config Config) (string, error) {
	claims := Claims{
		AccountID: accountID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Expires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

func ValidateToken(tokenString string, config Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %w", core_errors.ErrUnauthorized)
		}
		return []byte(config.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", core_errors.ErrUnauthorized)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims: %w", core_errors.ErrUnauthorized)
	}

	return claims, nil
}
