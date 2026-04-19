package domain

import "time"

type Task struct {
	ID      int
	Version int

	Title        string
	Description  *string
	Completed    bool
	Created_at   time.Time
	Completed_at *time.Time
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	created_at time.Time,
	completed_at *time.Time,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		Created_at:   created_at,
		Completed_at: completed_at,
	}
}

// func NewTaskUninitialized(title string, description *string, ) Task {
// 	return NewUser(
// 		UninitializedID,
// 		UninitializedVersion,
// 		fullName,
// 		phoneNumber,
// 	)
// }