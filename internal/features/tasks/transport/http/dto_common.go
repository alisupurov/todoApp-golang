package tasks_transport_http

import (
	"time"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID int        `json:"author_user_id"`
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

