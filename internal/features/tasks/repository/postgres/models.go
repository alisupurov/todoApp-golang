package tasks_postgres_repository

import (
	"time"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func TaskDomainsFromModels(taskModels []TaskModel) []domain.Task {
	domainTasks := make([]domain.Task, len(taskModels))
	for i, taskModel := range taskModels {
		domainTasks[i] = TaskDomainFromModel(taskModel)
	}
	return domainTasks
}

func TaskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}