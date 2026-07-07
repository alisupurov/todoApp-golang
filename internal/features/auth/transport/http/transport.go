package auth_transport_http

import (
	"context"
	"net/http"

	core_http_server "github.com/alisupurov/todoApp-golang/internal/core/transport/http/server"
)

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{authService: authService}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
	}
}
