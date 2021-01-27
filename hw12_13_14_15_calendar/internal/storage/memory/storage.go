package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventAlreadyExists = fmt.Errorf("event with such id already exists")
	ErrEventDoesNotExist  = fmt.Errorf("event with such id does not exist")
)

type Storage struct {
	mu     *sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:     new(sync.RWMutex),
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) Add(ctx context.Context, e storage.Event) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[e.ID]; ok {
		return e, ErrEventAlreadyExists
	}

	s.events[e.ID] = e

	return e, nil
}

func (s *Storage) Update(ctx context.Context, e storage.Event) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[e.ID]; !ok {
		return e, ErrEventDoesNotExist
	}

	s.events[e.ID] = e

	return e, nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return ErrEventDoesNotExist
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) ListPerDay(ctx context.Context, t time.Time) ([]storage.Event, error) {
	result := make([]storage.Event, 0)

	y, m, d := t.Date()

	beginningBorder := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	endOfBorder := time.Date(y, m, d, 23, 59, 59, int(time.Minute-time.Nanosecond), t.Location())

	for _, event := range s.events {
		if ok := checkDateBorders(event, beginningBorder, endOfBorder); ok {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) ListPerWeek(ctx context.Context, t time.Time) ([]storage.Event, error) {
	result := make([]storage.Event, 0)

	y, m, d := t.Date()

	beginningBorder := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	endOfBorder := beginningBorder.AddDate(0, 0, 7).Add(-time.Nanosecond)

	for _, event := range s.events {
		if ok := checkDateBorders(event, beginningBorder, endOfBorder); ok {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) ListPerMonth(ctx context.Context, t time.Time) ([]storage.Event, error) {
	result := make([]storage.Event, 0)

	y, m, d := t.Date()

	beginningBorder := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	endOfBorder := beginningBorder.AddDate(0, 1, 0).Add(-time.Nanosecond)

	for _, event := range s.events {
		if ok := checkDateBorders(event, beginningBorder, endOfBorder); ok {
			result = append(result, event)
		}
	}

	return result, nil
}

func checkDateBorders(event storage.Event, beginningOfDay time.Time, endOfDay time.Time) bool {
	isBeginningPassed := event.EventDate.After(beginningOfDay) || event.EventDate.Equal(endOfDay)
	isEndPassed := event.EventDate.Before(endOfDay) || (event.EventDate.Equal(endOfDay))

	return isBeginningPassed && isEndPassed
}
