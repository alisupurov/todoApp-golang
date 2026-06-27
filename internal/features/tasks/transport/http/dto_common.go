package tasks_transport_http

import (
	"time"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"             example:"5"`
	Version      int        `json:"version"        example:"1"`
	Title        string     `json:"title"          example:"Купить молоко"`
	Description  *string    `json:"description"    example:"2 литра"`
	Completed    bool       `json:"completed"      example:"false"`
	CreatedAt    time.Time  `json:"created_at"     example:"2026-06-21T21:33:15.393205Z"`
	CompletedAt  *time.Time `json:"completed_at"   example:"2026-06-21T21:34:00Z"`
	AuthorUserID int        `json:"author_user_id" example:"4"`
}

func taskDTOFromDomain(user domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           user.ID,
		Version:      user.Version,
		Title:        user.Title,
		Description:  user.Description,
		Completed:    user.Completed,
		CreatedAt:    user.CreatedAt,
		CompletedAt:  user.CompletedAt,
		AuthorUserID: user.AuthorUserId,
	}
}

func tasksDTOFromDomains(tasks []domain.Task) []TaskDTOResponse {
	tasksDTO := make([]TaskDTOResponse, len(tasks))
	for i, task := range tasks {
		tasksDTO[i] = taskDTOFromDomain(task)
	}

	return tasksDTO
}

