package tasks_transport_http

import (
	"net/http"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"              example:"Купить молоко"`
	Description  *string `json:"description" validate:"omitempty,min=3,max=1000"      example:"2 литра"`
	AuthorUserID int     `json:"author_user_id" validate:"required"                   example:"4"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask godoc
// @Summary      Создать задачу
// @Description  Создать новую задачу в системе
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body CreateTaskRequest true "CreateTask тело"
// @Success      201  {object}  CreateTaskResponse "Успешно созданная задача"
// @Failure      400  {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      404  {object}  core_http_response.ErrorResponse "Автор задачи не найден"
// @Failure      500  {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks [post]
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

