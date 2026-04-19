package users_transport_http

import (
	"net/http"

	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id path value")
		return
	}

	if err = h.usersService.DeleteUser(ctx, userId); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}
