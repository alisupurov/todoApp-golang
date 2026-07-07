package auth_transport_http

import (
	"net/http"

	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"        example:"ivan@example.com"`
	Password string `json:"password" validate:"required,min=6,max=72" example:"secret123"`
}

// Login godoc
// @Summary      Вход
// @Description  Войти в аккаунт, получить JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login тело"
// @Success      200  {object}  TokenResponse
// @Failure      400  {object}  core_http_response.ErrorResponse
// @Failure      401  {object}  core_http_response.ErrorResponse "Неверный email или пароль"
// @Failure      500  {object}  core_http_response.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	token, err := h.authService.Login(ctx, request.Email, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login")
		return
	}

	responseHandler.JSONResponse(TokenResponse{Token: token}, http.StatusOK)
}
