package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/lib/pq"
)

var (
	ErrEventIDIsExist    = errors.New("ivent ID is already exists")
	ErrEventIDIsNotExist = errors.New("ivent is not exists")
)

type StorageInDB struct {
	DB *sql.DB
}

// New returns a StorageInDB.
func New(ctx context.Context, connStr string, logg logger.Logger) (*StorageInDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	s := &StorageInDB{
		DB: db,
	}
	go s.Close(ctx, logg)
	return s, nil
}

// Close closes connection to DB.
func (s *StorageInDB) Close(ctx context.Context, logg logger.Logger) {
	<-ctx.Done()
	if err := s.DB.Close(); err != nil {
		logg.Error("failed to stop DB: " + err.Error())
	}
}

// CreateEvent adds event into storage.
func (s *StorageInDB) CreateEvent(ctx context.Context, event storage.Event) error {
	_, err := s.DB.ExecContext(ctx, "INSERT INTO events (ID, Title, DateTime, Duration, Description, UserID, NoteBefore) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		event.ID,
		event.Title,
		event.DateTime,
		event.Duration,
		event.Description,
		event.UserID,
		event.NoteBefore,
	)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return ErrEventIDIsExist
		default:
			return fmt.Errorf("creating of event has got an error: %w", err)
		}
	}
	return nil
}

// UpdateEvent updates event in storage.
func (s *StorageInDB) UpdateEvent(ctx context.Context, event storage.Event) error {
	err := s.isExistCheck(ctx, event)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx,
		"UPDATE events SET Title=$2, DateTime=$3, Duration=$4, Description=$5, UserID=$6, NoteBefore=$7	WHERE ID=$1",
		event.ID,
		event.Title,
		event.DateTime,
		event.Duration,
		event.Description,
		event.UserID,
		event.NoteBefore,
	)
	if err != nil {
		return fmt.Errorf("updating of event has got an error: %w", err)
	}
	return nil
}

// DeleteEvent deletes event from storage.
func (s *StorageInDB) DeleteEvent(ctx context.Context, event storage.Event) error {
	err := s.isExistCheck(ctx, event)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx,
		"DELETE FROM events WHERE ID=$1", event.ID)
	if err != nil {
		return fmt.Errorf("deleting of event has got an error: %w", err)
	}
	return nil
}

// isExistCheck checks wether event is exist.
func (s *StorageInDB) isExistCheck(ctx context.Context, event storage.Event) error {
	var eventInDB string
	row := s.DB.QueryRowContext(ctx, "SELECT ID FROM events WHERE ID=$1", event.ID)
	if errors.Is(row.Scan(&eventInDB), sql.ErrNoRows) {
		return ErrEventIDIsNotExist
	}
	return nil
}

// GetEventsForDay return slice of events for one day.
func (s *StorageInDB) GetEventsForDay(ctx context.Context, dateOfQuestion time.Time) ([]storage.Event, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT * FROM events WHERE EXTRACT(YEAR FROM datetime)=$1 AND EXTRACT(MONTH FROM datetime)=$2 AND EXTRACT(DAY FROM datetime)=$3", dateOfQuestion.Year(), dateOfQuestion.Month(), dateOfQuestion.Day())
	if err != nil {
		return nil, fmt.Errorf("getting events for the day has got an error: %w", err)
	}
	defer rows.Close()
	return postWorkFor(rows)
}

// GetEventsForWeek returns slice of events for one week, you need to determine any day from that week.
func (s *StorageInDB) GetEventsForWeek(ctx context.Context, dateOfQuestion time.Time) ([]storage.Event, error) {
	year, week := dateOfQuestion.ISOWeek()
	rows, err := s.DB.QueryContext(ctx, "SELECT * FROM events WHERE EXTRACT(YEAR FROM datetime)=$1 AND EXTRACT(WEEK FROM datetime)=$2", year, week)
	if err != nil {
		return nil, fmt.Errorf("getting events for the week has got an error: %w", err)
	}
	defer rows.Close()
	return postWorkFor(rows)
}

// GetEventsForMonth returns slice of events for one month, you need to determine any day from that month.
func (s *StorageInDB) GetEventsForMonth(ctx context.Context, dateOfQuestion time.Time) ([]storage.Event, error) {
	year, month, _ := dateOfQuestion.Date()
	rows, err := s.DB.QueryContext(ctx, "SELECT * FROM events WHERE EXTRACT(YEAR FROM datetime)=$1 AND EXTRACT(MONTH FROM datetime)=$2", year, month)
	if err != nil {
		return nil, fmt.Errorf("getting events for the month has got an error: %w", err)
	}
	defer rows.Close()
	return postWorkFor(rows)
}

// postWorkFor performs a slice of events from sql.Rows.
func postWorkFor(rows *sql.Rows) ([]storage.Event, error) {
	var sliceOfEvents []storage.Event
	for rows.Next() {
		var e storage.Event
		rows.Scan(&e.ID, &e.Title, &e.DateTime, &e.Duration, &e.Description, &e.UserID, &e.NoteBefore)
		sliceOfEvents = append(sliceOfEvents, e)
	}
	err := rows.Err()
	return sliceOfEvents, err
}
