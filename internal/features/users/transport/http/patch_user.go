package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
	core_http_types "github.com/alisupurov/todoApp-golang/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"    example:"Иван Иванов"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" example:"+79995553322"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("full_name cannot be null: %w", core_errors.ErrInvalidArgument)
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`Full Name` must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set && r.PhoneNumber.Value != nil {
		phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("`Phone Number` must be between 10 and 15 symbols")
		}

		if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
			return fmt.Errorf("`Phone Number` must starts with `+`")
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser godoc
// @Summary      Обновить пользователя
// @Description  Частично обновить данные пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Param        request body PatchUserRequest true "PatchUser тело"
// @Success      200  {object}  PatchUserResponse "Обновлённый пользователь"
// @Failure      400  {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      404  {object}  core_http_response.ErrorResponse "Not found"
// @Failure      409  {object}  core_http_response.ErrorResponse "Conflict"
// @Failure      500  {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id path value")
		return
	}


	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := patchUserDomainFromRequest(request)
	domain, err := h.usersService.PatchUser(ctx, id, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(domain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func patchUserDomainFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(), 
		request.PhoneNumber.ToDomain(),
	)
}