package sqlstorage

import (
	"context"
	"fmt"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, cfg config.Config) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", cfg.Storage.SQL.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) generateID() uuid.UUID {
	return uuid.New()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (uuid.UUID, error) {
	eventUuid := s.generateID()
	result, err := s.db.ExecContext(ctx, "INSERT INTO events (id, title) VALUES ($1, $2)", eventUuid.String(), event.Title)
	if err != nil {
		return uuid.UUID{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return uuid.UUID{}, err
	}

	if rows != 1 {
		return uuid.UUID{}, fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return eventUuid, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) error {
	result, err := s.db.ExecContext(ctx, "UPDATE events SET title=$1 WHERE id=$2", event.Title, id.String())
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id=$1", id.String())
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (s *Storage) GetEvents(ctx context.Context, date time.Time, offset int) ([]storage.Event, error) {
	startDate := date.Format(time.DateOnly)
	endDate := date.AddDate(0, 0, offset).Format(time.DateOnly)

	result := s.db.QueryRowxContext(ctx, "SELECT id, title, start_time, description FROM events WHERE start_time BETWEEN $1 AND $2", startDate, endDate)
	defer s.db.Close()

	var events []storage.Event
	err := result.StructScan(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}
