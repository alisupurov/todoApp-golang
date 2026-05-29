package tasks_transport_http

import (
	"net/http"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID int     `json:"author_user_id" validate:"required=true"`
}

type CreateTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	domainTask := domain.NewTaskUninitialized(request.Title, request.Description, request.AuthorUserID)

	domainTask, err := h.tasksService.CreateTask(ctx, domainTask)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)

		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(domainTask))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

