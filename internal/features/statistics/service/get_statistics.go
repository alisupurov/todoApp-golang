package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_errors "github.com/alisupurov/todoApp-golang/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"'to' must be after 'from: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userId, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("failed to get statistics: %w", err)
	}

	statistics := calcStatisctics(tasks)

	return statistics, nil
}	

func calcStatisctics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)
	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDiration := task.CompletionDiration()
		if completionDiration != nil {
			totalCompletionDuration += *completionDiration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageComplionTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)
		tasksAverageComplionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedRate,
		tasksAverageComplionTime,
	)
}