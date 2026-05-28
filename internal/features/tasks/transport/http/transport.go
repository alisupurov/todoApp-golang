package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_http_server "github.com/alisupurov/todoApp-golang/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTask(
		ctx context.Context,
		id int,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		limit *int,
		offset *int,
		user_id *int,
	) ([]domain.Task, error)

	PatchTask(
		ctx context.Context,
		id int,
		taskPatch domain.TaskPatch,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id int,
	) error
}

func NewTasksHTTPHandler(
	tasksService TasksService,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method: http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
		{
			Method: http.MethodPatch,
			Path: "task/{id}",
			Handler: h.PatchTask,
		},
	}
}
