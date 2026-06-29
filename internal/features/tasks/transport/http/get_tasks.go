package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary      Получить список задач
// @Description  Получить список задач с пагинацией и фильтром по автору
// @Tags         tasks
// @Produce      json
// @Param        limit query int false "Лимит записей"
// @Param        offset query int false "Смещение"
// @Param        user_id query int false "ID автора задачи"
// @Success      200  {array}   TaskDTOResponse "Список задач"
// @Failure      400  {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      500  {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, userId, err := getUserIdLimitOffsetFromRequest(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit, offset and user_id")
		return
	}

	domainTasks, err := h.tasksService.GetTasks(ctx, limit, offset, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := GetTasksResponse(tasksDTOFromDomains(domainTasks))
	responseHandler.JSONResponse(response, http.StatusOK)
}


func getUserIdLimitOffsetFromRequest(r *http.Request) (*int, *int, *int, error) {
	const (
		userIdQueryParam = "user_id"
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get limit query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get offset query param: %w", err)
	}

	userId, err := core_http_request.GetIntQueryParam(r, userIdQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get user_id query param: %w", err)
	}

	return limit, offset, userId, nil
}
