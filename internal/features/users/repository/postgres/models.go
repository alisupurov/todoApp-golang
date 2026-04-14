package users_postgres_repository

import "github.com/alisupurov/todoApp-golang/internal/core/domain"

type UserModel struct {
	ID           int
	Version      int
	Full_name    string
	Phone_number *string
}

func userDomainsFromModels(userModels []UserModel) []domain.User {
	users := make([]domain.User, len(userModels))
	for i, user := range userModels {
		users[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.Full_name,
			user.Phone_number,
		)
	}
	return users
}
