package sqlstorage

import (
	"context"
	"fmt"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strings"
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

func (s *Storage) CreateEvent(ctx context.Context, event model.Event) (uuid.UUID, error) {
	eventUuid := s.generateID()

	result, err := s.db.ExecContext(
		ctx,
		`INSERT INTO event (id, title, start_time, description, duration, notify_before, user_id) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		eventUuid.String(),
		event.Title,
		event.StartTime,
		event.Description,
		event.Duration,
		event.NotifyBefore,
		event.UserID,
	)

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

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, event model.Event) error {
	result, err := s.db.ExecContext(ctx, "UPDATE event SET title=$1 WHERE id=$2", event.Title, id.String())
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
	result, err := s.db.ExecContext(ctx, "DELETE FROM event WHERE id=$1", id.String())
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

func (s *Storage) GetEvents(ctx context.Context, date time.Time, offset int) ([]model.Event, error) {
	startDate := date.Format(time.DateOnly)
	endDate := date.AddDate(0, 0, offset).Format(time.DateOnly)

	result := s.db.QueryRowxContext(ctx, "SELECT id, title, start_time, description, duration, notify_before, user_id FROM event WHERE start_time BETWEEN $1 AND $2", startDate, endDate)
	defer s.db.Close()

	var events []model.Event
	err := result.StructScan(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) GetNotifications(ctx context.Context, date time.Time) ([]model.Notification, error) {
	dateString := date.Format("2006-01-02 15:04:05")

	result := s.db.QueryRowxContext(ctx, `
		SELECT id, title, start_time, user_id FROM event WHERE notify_before IS NOT NULL AND sent = false and start_time - notify_before <= $1 AND start_time > $1
	`,
		dateString,
	)
	defer s.db.Close()

	var notifications []model.Notification
	err := result.StructScan(&notifications)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *Storage) MarkEventsAsNotified(ctx context.Context, notifications []model.Notification) error {
	eventIDs := make([]string, 0, len(notifications))

	for _, notification := range notifications {
		eventIDs = append(eventIDs, "'"+notification.EventID.String()+"'")
	}

	eventIDsString := strings.Join(eventIDs, ",")

	result, err := s.db.ExecContext(ctx, "UPDATE event SET sent=true WHERE id=$1", eventIDsString)
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

func (s *Storage) DeleteOldEvents(ctx context.Context) error {
	deleteAfter := time.Now().AddDate(-1, 0, 0)

	result, err := s.db.ExecContext(ctx, "DELETE FROM event WHERE start_time < $1", deleteAfter)
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
