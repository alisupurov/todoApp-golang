package domain

import "time"

type Statistics struct {
	TasksCreated             int
	TasksCompleted           int
	TaskCompletedRate        *float64
	TasksAverageCompleteTime *time.Duration
}

func NewStatistics(
	tasksCreated int,
	tasksCompleted int,
	taskCompletedRate *float64,
	tasksAverageCompleteTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated: tasksCreated,
		TasksCompleted: tasksCompleted,
		TaskCompletedRate: taskCompletedRate,
		TasksAverageCompleteTime: tasksAverageCompleteTime,
	}
}