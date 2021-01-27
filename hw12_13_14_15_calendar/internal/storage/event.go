package storage

import (
	"time"
)

type Event struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	EventDate   time.Time     `json:"eventDate"`
	Duration    time.Duration `json:"duration"`
	Description string        `json:"description"`
	UserID      string        `json:"userId"`
	RemindIn    time.Duration `json:"remindIn"`
}
