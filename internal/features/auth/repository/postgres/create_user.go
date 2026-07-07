package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	core_postgres_pool "github.com/alisupurov/todoApp-golang/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) CreateAccount(ctx context.Context, email, passwordHash string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO todoapp.accounts (email, password_hash)
		VALUES ($1, $2)
		RETURNING id;
	`

	var accountID int

	row := r.pool.QueryRow(ctx, query, email, passwordHash)
	if err := row.Scan(&accountID); err != nil {
		if errors.Is(err, core_postgres_pool.ErrUniqueViolation) {
			return 0, fmt.Errorf("account with email=%s already exists: %w", email, core_errors.ErrConflict)
		}
		return 0, fmt.Errorf("scan error: %w", err)
	}

	return accountID, nil
}
