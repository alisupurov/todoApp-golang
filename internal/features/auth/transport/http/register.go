package auth_transport_http

import (
	"net/http"

	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"        example:"ivan@example.com"`
	Password string `json:"password" validate:"required,min=6,max=72" example:"secret123"`
}

type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiJ9..."`
}

// Register godoc
// @Summary      Регистрация
// @Description  Создать аккаунт для входа в систему, получить JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Register тело"
// @Success      201  {object}  TokenResponse
// @Failure      400  {object}  core_http_response.ErrorResponse
// @Failure      409  {object}  core_http_response.ErrorResponse "Email уже занят"
// @Failure      500  {object}  core_http_response.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHTTPHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request RegisterRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	token, err := h.authService.Register(ctx, request.Email, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to register")
		return
	}

	responseHandler.JSONResponse(TokenResponse{Token: token}, http.StatusCreated)
}
