package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeoutPgx())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	ORDER BY id
	LIMIT $1 
	OFFSET $2;`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []UserModel
	for rows.Next() {
		var userModel UserModel
		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.Full_name,
			&userModel.Phone_number,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, userModel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over users: %w", err)
	}

	userDomains := userDomainsFromModels(users)
	return userDomains, nil
}