package model

import "time"

type Task struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
