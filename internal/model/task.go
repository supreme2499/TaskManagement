package model

import "time"

type Task struct {
	ID          int
	NameTask    string
	Description string
	Status      string
	Deadline    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
