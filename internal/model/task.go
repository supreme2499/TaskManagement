package model

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	Deadline    time.Time
	Created_at  time.Time
	Updated_at  time.Time
}
