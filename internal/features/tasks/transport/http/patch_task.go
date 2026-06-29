package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_http_request "github.com/alisupurov/todoApp-golang/internal/core/transport/http/request"
	core_http_response "github.com/alisupurov/todoApp-golang/internal/core/transport/http/response"
	core_http_types "github.com/alisupurov/todoApp-golang/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"       example:"Новое название"`
	Description core_http_types.Nullable[string] `json:"description" example:"Новое описание"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"   example:"true"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("full_name cannot be null: %w", core_errors.ErrInvalidArgument)
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set && r.Description.Value != nil {
		descriptionLen := len([]rune(*r.Description.Value))
		if descriptionLen < 1 && descriptionLen > 1000 {
			return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
		}
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

// PatchTask godoc
// @Summary      Обновить задачу
// @Description  Частично обновить данные задачи
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path int true "ID задачи"
// @Param        request body PatchTaskRequest true "PatchTask тело"
// @Success      200  {object}  PatchTaskResponse "Обновлённая задача"
// @Failure      400  {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      404  {object}  core_http_response.ErrorResponse "Not found"
// @Failure      409  {object}  core_http_response.ErrorResponse "Conflict"
// @Failure      500  {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id in path value")
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskPatch := patchTaskDomainFromRequest(request)
	domainTask, err := h.tasksService.PatchTask(ctx, id, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(domainTask))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func patchTaskDomainFromRequest(request PatchTaskRequest) domain.TaskPatch{
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
