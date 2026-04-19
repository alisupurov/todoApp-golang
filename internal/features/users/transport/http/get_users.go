package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Debug("invoce GetUsers handler")

	limit, offset, err := getLimitOffsetFromRequest(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit and offset query params")
		return
	}

	users, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(usersDTOFromDomain(users))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetFromRequest(r *http.Request) (limit, offset *int, err error) {
	const (
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)

	limit, err = core_http_request.GetIntQueryParam(r, limitQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get limit query param: %w", err)
	}

	offset, err = core_http_request.GetIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get offset query param: %w", err)
	}

	return limit, offset, nil
}
