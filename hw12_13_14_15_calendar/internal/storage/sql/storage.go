package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/utils/idgen"
	_ "github.com/jackc/pgx/stdlib" // nolint: gci
)

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error while exacute sql.Open: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Connect(ctx context.Context, dsn string) (err error) {
	err = s.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("error while connect to the db: %w", err)
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) Add(ctx context.Context, e storage.Event) (storage.Event, error) {
	query := "INSERT INTO events (id, title, event_date, duration, description, user_id, remind_in) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"

	e.ID, _ = idgen.PrefixedID("EV")

	_, err := s.db.ExecContext(
		ctx,
		query,
		e.ID,
		e.Title,
		e.EventDate,
		int(e.Duration),
		e.Description,
		e.UserID,
		int(e.RemindIn),
	)

	if err != nil {
		err = fmt.Errorf("error while adding event: %w", err)
	}

	return e, err
}

func (s *Storage) Update(ctx context.Context, e storage.Event) (storage.Event, error) {
	query := `UPDATE events SET 
                  title=$1, 
                  event_date=$2, 
                  duration=$3, 
                  description=$4, 
                  user_id=$5, 
                  remind_in=$6
				WHERE id=$7;`

	_, err := s.db.ExecContext(
		ctx,
		query,
		e.Title,
		e.EventDate,
		e.Duration,
		e.Description,
		e.UserID,
		e.RemindIn,
		e.ID,
	)

	if err != nil {
		err = fmt.Errorf("error while updating event: %w", err)
	}

	return e, err
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id=$1;`

	result, err := s.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		err = fmt.Errorf("error while deleting event: %w", err)
	} else {
		ra, _ := result.RowsAffected()
		if ra == 0 {
			err = fmt.Errorf("error while deleting event: no row with such id")
		}
	}

	return err
}

func (s *Storage) ListPerDay(ctx context.Context, t time.Time) ([]storage.Event, error) {
	y, m, d := t.Date()

	beginningBorder := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	endOfBorder := time.Date(y, m, d, 23, 59, 59, int(time.Minute-time.Nanosecond), t.Location())

	result, err := s.listPerPeriod(ctx, beginningBorder, endOfBorder)
	if err != nil {
		return nil, fmt.Errorf("error on event list request: %w", err)
	}

	return result, nil
}

func (s *Storage) ListPerWeek(ctx context.Context, t time.Time) ([]storage.Event, error) {
	y, m, d := t.Date()

	beginningBorder := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	endOfBorder := beginningBorder.AddDate(0, 0, 7).Add(-time.Nanosecond)

	result, err := s.listPerPeriod(ctx, beginningBorder, endOfBorder)
	if err != nil {
		return nil, fmt.Errorf("error on event list request: %w", err)
	}

	return result, nil
}

func (s *Storage) ListPerMonth(ctx context.Context, t time.Time) ([]storage.Event, error) {
	y, m, d := t.Date()

	beginningBorder := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	endOfBorder := beginningBorder.AddDate(0, 1, 0).Add(-time.Nanosecond)

	result, err := s.listPerPeriod(ctx, beginningBorder, endOfBorder)
	if err != nil {
		return nil, fmt.Errorf("error on event list request: %w", err)
	}

	return result, nil
}

func (s *Storage) listPerPeriod(ctx context.Context, start time.Time, end time.Time) ([]storage.Event, error) {
	query := `SELECT * FROM events WHERE event_date>=$1 and event_date<=$2;`

	rows, err := s.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("error while exec query: %w", err)
	}
	defer rows.Close()

	result := make([]storage.Event, 0)
	for rows.Next() {
		ev := storage.Event{}
		err = rows.Scan(&ev.ID, &ev.Title, &ev.EventDate, &ev.Duration, &ev.Description, &ev.UserID, &ev.RemindIn)
		if err != nil {
			return nil, fmt.Errorf("error while scanning event rows: %w", err)
		}
		result = append(result, ev)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error while getting events: %w", err)
	}

	return result, nil
}
