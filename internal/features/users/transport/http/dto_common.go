package users_transport_http

import "github.com/alisupurov/todoApp-golang/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id"            example:"10"`
	Version     int     `json:"version"       example:"3"`
	FullName    string  `json:"full_name"     example:"Иван Иванов"`
	PhoneNumber *string `json:"phone_number"  example:"79995553322"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomain(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))
	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}
	return usersDTO
}