package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeoutPgx())
	defer cancel()

	query := `
	UPDATE todoapp.users
	SET 
		full_name = $1, 
		phone_number = $2, 
		version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber, id, user.Version)

	var userModel UserModel	
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Full_name,
		&userModel.Phone_number,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id=%d concurrently acessed: %w", id, core_errors.ErrConflict)
		} else {
			return domain.User{}, fmt.Errorf("scan error")
		}
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Full_name,
		userModel.Phone_number,
	)
	
	return userDomain, nil
}