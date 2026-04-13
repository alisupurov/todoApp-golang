package users_postgres_repository

type UserModel struct {
	ID           int
	Version      int
	Full_name    string
	Phone_number *string
}
