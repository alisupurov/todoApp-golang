package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	limit *int,
	offset *int,
	user_id *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	%s 
	ORDER BY id ASC
	LIMIT $1 
	OFFSET $2;
	`

	args := []any{limit, offset}

	if user_id != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id = $3")
		args = append(args, user_id)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	var tasks []TaskModel

	for rows.Next() {
		var taskModel TaskModel
		err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		tasks = append(tasks, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := TaskDomainsFromModels(tasks)
	return taskDomains, nil
}