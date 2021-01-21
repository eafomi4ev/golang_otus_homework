package storage

import (
	"time"
)

type Event struct {
	ID          string
	Title       string
	EventDate   time.Time
	Duration    time.Duration
	Description string
	UserID      string
	RemindIn    time.Duration
}
