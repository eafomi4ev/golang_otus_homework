package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/utils/idgen"
)

func TestStorage(t *testing.T) {
	t.Run("Add", func(t *testing.T) {

		eventId, err := idgen.PrefixedID("EV")
		require.NoError(t, err)
		require.Len(t, eventId, 15)

		userID, err := idgen.PrefixedID("USR")
		require.NoError(t, err)
		require.Len(t, userID, 16)

		e := storage.Event{
			ID:          eventId,
			Title:       "Test Event",
			EventDate:   time.Now(),
			Duration:    time.Hour,
			Description: "Event just for test",
			UserID:      userID,
			RemindIn:    0,
		}

		s := New()

		ctx := context.Background()

		require.Len(t, s.events, 0)

		err = s.Add(ctx, e)

		require.NoError(t, err)
		require.Len(t, s.events, 1)
	})

	t.Run("Update", func(t *testing.T) {

		eventId, err := idgen.PrefixedID("EV")
		require.NoError(t, err)
		require.Len(t, eventId, 15)

		userID, err := idgen.PrefixedID("USR")
		require.NoError(t, err)
		require.Len(t, userID, 16)

		e := storage.Event{
			ID:          eventId,
			Title:       "Test Event",
			EventDate:   time.Now(),
			Duration:    time.Hour,
			Description: "Event just for test",
			UserID:      userID,
			RemindIn:    0,
		}
		s := New()
		require.Len(t, s.events, 0)

		ctx := context.Background()
		err = s.Add(ctx, e)
		require.NoError(t, err)
		require.Equal(t, s.events[e.ID].Title, "Test Event")

		e.Title = "Updated Title event"
		err = s.Update(ctx, e)
		require.NoError(t, err)
		require.Equal(t, s.events[e.ID].Title, "Updated Title event")
	})

	t.Run("Delete", func(t *testing.T) {

		eventId, err := idgen.PrefixedID("EV")
		require.NoError(t, err)

		userID, err := idgen.PrefixedID("USR")
		require.NoError(t, err)

		e := storage.Event{
			ID:          eventId,
			Title:       "Test Event",
			EventDate:   time.Now(),
			Duration:    time.Hour,
			Description: "Event just for test",
			UserID:      userID,
			RemindIn:    0,
		}
		s := New()
		require.Len(t, s.events, 0)

		ctx := context.Background()
		err = s.Add(ctx, e)
		require.NoError(t, err)
		require.Len(t, s.events, 1)

		err = s.Delete(ctx, e.ID)
		require.NoError(t, err)
		require.Len(t, s.events, 0)
	})

	t.Run("List per day", func(t *testing.T) {

		eventId1, err := idgen.PrefixedID("EV")
		require.NoError(t, err)
		eventId2, err := idgen.PrefixedID("EV")
		require.NoError(t, err)

		userID, err := idgen.PrefixedID("USR")
		require.NoError(t, err)

		e1 := storage.Event{
			ID:          eventId1,
			Title:       "Test Event 1",
			EventDate:   time.Now().AddDate(0, 0, -3),
			Duration:    time.Hour,
			Description: "Event just for test",
			UserID:      userID,
			RemindIn:    0,
		}

		e2 := storage.Event{
			ID:          eventId2,
			Title:       "Test Event 2",
			EventDate:   time.Now().AddDate(0, -1, 0),
			Duration:    time.Hour,
			Description: "Event just for test",
			UserID:      userID,
			RemindIn:    0,
		}
		s := New()
		require.Len(t, s.events, 0)

		ctx := context.Background()
		err = s.Add(ctx, e1)
		require.NoError(t, err)
		err = s.Add(ctx, e2)
		require.NoError(t, err)

		require.Len(t, s.events, 2)

		events, err := s.ListPerDay(ctx, time.Now().AddDate(0, 0, -3))
		require.Len(t, events, 1)

		events, err = s.ListPerDay(ctx, time.Now().AddDate(0, 0, -4))
		require.Len(t, events, 0)

		events, err = s.ListPerDay(ctx, time.Now().AddDate(0, 0, -2))
		require.Len(t, events, 0)
	})

	t.Run("List per week", func(t *testing.T) {

		eventId1, err := idgen.PrefixedID("EV")
		require.NoError(t, err)
		eventId2, err := idgen.PrefixedID("EV")
		require.NoError(t, err)

		userID, err := idgen.PrefixedID("USR")
		require.NoError(t, err)

		e1 := storage.Event{
			ID:          eventId1,
			Title:       "Test Event 1",
			EventDate:   time.Now(),
			Duration:    time.Hour,
			Description: "Event just for test",
			UserID:      userID,
			RemindIn:    0,
		}

		e2 := storage.Event{
			ID:          eventId2,
			Title:       "Test Event 2",
			EventDate:   time.Now().AddDate(0, -1, 0),
			Duration:    time.Hour,
			Description: "Event just for test",
			UserID:      userID,
			RemindIn:    0,
		}
		s := New()
		require.Len(t, s.events, 0)

		ctx := context.Background()
		err = s.Add(ctx, e1)
		require.NoError(t, err)
		err = s.Add(ctx, e2)
		require.NoError(t, err)

		require.Len(t, s.events, 2)

		events, err := s.ListPerWeek(ctx, time.Now().AddDate(0, 0, -6))
		require.Len(t, events, 1)

		events, err = s.ListPerWeek(ctx, time.Now().AddDate(0, 0, -7))
		require.Len(t, events, 0)
	})

	t.Run("List per month", func(t *testing.T) {

		eventId1, err := idgen.PrefixedID("EV")
		require.NoError(t, err)
		eventId2, err := idgen.PrefixedID("EV")
		require.NoError(t, err)

		userID, err := idgen.PrefixedID("USR")
		require.NoError(t, err)

		d1, _ := time.Parse("2006-01-02", "2020-01-31")
		e1 := storage.Event{
			ID:          eventId1,
			Title:       "Test Event 1",
			EventDate:   d1,
			Duration:    time.Hour,
			Description: "Event just for test",
			UserID:      userID,
			RemindIn:    0,
		}

		d2, _ := time.Parse("2006-01-02", "2020-02-15")
		e2 := storage.Event{
			ID:          eventId2,
			Title:       "Test Event 2",
			EventDate:   d2,
			Duration:    time.Hour,
			Description: "Event just for test",
			UserID:      userID,
			RemindIn:    0,
		}
		s := New()
		require.Len(t, s.events, 0)

		ctx := context.Background()
		err = s.Add(ctx, e1)
		require.NoError(t, err)
		err = s.Add(ctx, e2)
		require.NoError(t, err)

		require.Len(t, s.events, 2)

		monthStartDate, _ := time.Parse("2006-01-02", "2020-01-01")
		events, err := s.ListPerMonth(ctx, monthStartDate)
		require.Len(t, events, 1)

		events, err = s.ListPerMonth(ctx, monthStartDate.AddDate(0, 0, -1))
		require.Len(t, events, 0)
	})
}
