package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_postgres_pool "github.com/alisupurov/todoApp-golang/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeoutPgx())
	defer cancel()

	query := `
		SELECT id, version, full_name, phone_number
		FROM todoapp.users
		WHERE id = $1;
	`
	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Full_name,
		&userModel.Phone_number,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id %d not found: %w", id, core_postgres_pool.ErrNoRows)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Full_name,
		userModel.Phone_number,
	)
	return userDomain, nil
}
