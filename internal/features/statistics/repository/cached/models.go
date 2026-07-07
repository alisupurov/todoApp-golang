package statistics_cached_repository

import (
	"encoding/json"
	"time"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
)

type TaskModel struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID int        `json:"author_user_id"`
}

type TaskListModel []TaskModel

func (m TaskListModel) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *TaskListModel) Deserialize(data []byte) error {
	return json.Unmarshal(data, m)
}

func domainsToModel(tasks []domain.Task) TaskListModel {
	model := make(TaskListModel, len(tasks))
	for i, task := range tasks {
		model[i] = TaskModel{
			ID:           task.ID,
			Version:      task.Version,
			Title:        task.Title,
			Description:  task.Description,
			Completed:    task.Completed,
			CreatedAt:    task.CreatedAt,
			CompletedAt:  task.CompletedAt,
			AuthorUserID: task.AuthorUserId,
		}
	}

	return model
}

func modelToDomains(model TaskListModel) []domain.Task {
	tasks := make([]domain.Task, len(model))
	for i, taskModel := range model {
		tasks[i] = domain.NewTask(
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

	return tasks
}
