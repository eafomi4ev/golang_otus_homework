package app

import (
	"context"
	"fmt"
	"time"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/utils/idgen"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	Add(ctx context.Context, e storage.Event) (storage.Event, error)
	Update(ctx context.Context, e storage.Event) (storage.Event, error)
	Delete(ctx context.Context, id string) error
	ListPerDay(ctx context.Context, t time.Time) ([]storage.Event, error)
	ListPerWeek(ctx context.Context, day time.Time) ([]storage.Event, error)
	ListPerMonth(ctx context.Context, day time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, title string) (*storage.Event, error) {
	id, err := idgen.PrefixedID("EV")
	if err != nil {
		return nil, fmt.Errorf("cannot create event: %w", err)
	}

	userID, err := idgen.PrefixedID("USR")
	if err != nil {
		return nil, fmt.Errorf("cannot create user id: %w", err)
	}

	event := &storage.Event{
		ID:          id,
		Title:       title,
		EventDate:   time.Time{},
		Duration:    0,
		Description: "",
		UserID:      userID,
		RemindIn:    0,
	}

	return event, nil
}
