package tasks_service

import (
	"context"
	"fmt"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	id int,
	taskPatch domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	if err := task.ApplyPatch(taskPatch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return patchedTask, nil
}