package model

import "time"

type NotificationMessage struct {
	Event        string
	Timestamp    time.Time
	TaskID       int
	UserID       int
	ChangeStatus string
}
