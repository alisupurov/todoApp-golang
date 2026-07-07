package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	core_postgres_pool "github.com/alisupurov/todoApp-golang/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) GetAccountByEmail(ctx context.Context, email string) (int, string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, password_hash
		FROM todoapp.accounts
		WHERE email = $1;
	`

	var accountID int
	var passwordHash string

	row := r.pool.QueryRow(ctx, query, email)
	if err := row.Scan(&accountID, &passwordHash); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return 0, "", fmt.Errorf("account with email=%s: %w", email, core_errors.ErrNotFound)
		}
		return 0, "", fmt.Errorf("scan error: %w", err)
	}

	return accountID, passwordHash, nil
}
