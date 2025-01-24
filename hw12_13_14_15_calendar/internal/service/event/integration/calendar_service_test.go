package integration

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/service/event"
	sqlstorage "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type EventService interface {
	CreateEvent(ctx context.Context, event model.Event) (uuid.UUID, error)
	DayEventList(ctx context.Context, date time.Time) ([]model.Event, error)
	WeekEventList(ctx context.Context, date time.Time) ([]model.Event, error)
	MonthEventList(ctx context.Context, date time.Time) ([]model.Event, error)
}

type CalendarServiceSuite struct {
	suite.Suite
	svc EventService
	db  *sqlx.DB
}

func TestServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(CalendarServiceSuite))
}

func (s *CalendarServiceSuite) SetupSuite() {
	db, err := sqlx.ConnectContext(context.Background(), "postgres", "postgres://test:test@pg:5432/calendar?sslmode=disable")
	s.Require().NoError(err)

	repo := sqlstorage.New(db)
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	s.svc = event.NewEventService(*logger, repo)
	s.db = db
}

func (s *CalendarServiceSuite) SetupTest() {
	_, err := s.db.ExecContext(context.Background(), "TRUNCATE TABLE event")
	s.Require().NoError(err)
}

func (s *CalendarServiceSuite) TearDownSuite() {
	_, _ = s.db.ExecContext(context.Background(), "TRUNCATE TABLE event")
	s.db.Close()
}

func (s *CalendarServiceSuite) TestCreateEvent() {
	const (
		eventTitle = "test create event 1"
	)

	startTime := time.Now()
	m := model.Event{
		Title:     eventTitle,
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "user1",
	}

	eventID, err := s.svc.CreateEvent(context.Background(), m)
	s.Require().NoError(err)
	s.Require().NotEmpty(eventID)

	createdEvent := s.getDirectItem(m.Title)

	s.Require().NoError(err)
	s.Require().NotEmpty(createdEvent)
	s.Require().Equal(eventID, createdEvent.ID)
	s.Require().Equal(m.Title, createdEvent.Title)
	s.Require().Equal(m.StartTime.Format(time.DateTime), createdEvent.StartTime.Format(time.DateTime))
	s.Require().Equal(m.Duration, createdEvent.Duration)
	s.Require().Equal(m.UserID, createdEvent.UserID)
}

func (s *CalendarServiceSuite) TestDayEventList() {
	startTime := time.Now().Truncate(24 * time.Hour).Add(9 * time.Hour)
	m1 := model.Event{
		Title:     "event 1",
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m2 := model.Event{
		Title:     "event 2",
		StartTime: startTime.Add(time.Hour * 2),
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m3 := model.Event{
		Title:     "event 3",
		StartTime: startTime.Add(time.Hour * 25),
		Duration:  time.Hour,
		UserID:    "user1",
	}
	s.createDirectItem(m1)
	s.createDirectItem(m2)
	s.createDirectItem(m3)

	dayEvents, err := s.svc.DayEventList(context.Background(), time.Now())
	s.Require().NoError(err)
	s.Require().Equal(2, len(dayEvents))
}

func (s *CalendarServiceSuite) TestWeekEventList() {
	startTime := time.Now().Truncate(24 * time.Hour).Add(9 * time.Hour)
	m1 := model.Event{
		Title:     "event 1",
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m2 := model.Event{
		Title:     "event 2",
		StartTime: startTime.Add(time.Hour * 2),
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m3 := model.Event{
		Title:     "event 3",
		StartTime: startTime.Add(time.Hour * 25),
		Duration:  time.Hour,
		UserID:    "user1",
	}
	s.createDirectItem(m1)
	s.createDirectItem(m2)
	s.createDirectItem(m3)

	dayEvents, err := s.svc.WeekEventList(context.Background(), time.Now())
	s.Require().NoError(err)
	s.Require().Equal(3, len(dayEvents))
}

func (s *CalendarServiceSuite) TestMonthEventList() {
	startTime := time.Now().Truncate(24 * time.Hour).Add(9 * time.Hour)
	m1 := model.Event{
		Title:     "event 1",
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m2 := model.Event{
		Title:     "event 2",
		StartTime: startTime.Add(time.Hour * 25),
		Duration:  time.Hour,
		UserID:    "user1",
	}

	m3 := model.Event{
		Title:     "event 3",
		StartTime: startTime.Add(time.Hour * 24 * 30),
		Duration:  time.Hour,
		UserID:    "user1",
	}
	s.createDirectItem(m1)
	s.createDirectItem(m2)
	s.createDirectItem(m3)

	dayEvents, err := s.svc.MonthEventList(context.Background(), time.Now())
	s.Require().NoError(err)
	s.Require().Equal(3, len(dayEvents))
}
func (s *CalendarServiceSuite) getDirectItem(title string) model.Event {
	row := s.db.QueryRowxContext(context.Background(), "SELECT id, title, start_time, description, duration, user_id FROM event WHERE title = $1", title)

	var e model.Event
	var durationStr string

	err := row.Scan(
		&e.ID,
		&e.Title,
		&e.StartTime,
		&e.Description,
		&durationStr,
		&e.UserID,
	)

	if err != nil {
		s.Fail(err.Error())
	}

	duration, err := parseDuration(durationStr)
	if err != nil {
		s.Fail(err.Error())
	}
	e.Duration = duration

	return e
}

func (s *CalendarServiceSuite) createDirectItem(event model.Event) {
	eventId := uuid.New()

	result, err := s.db.ExecContext(
		context.Background(),
		`INSERT INTO event (id, title, start_time, description, duration, notify_before, user_id) 
			VALUES ($1, $2, $3, $4, $5::interval, $6::interval, $7)`,
		eventId.String(),
		event.Title,
		event.StartTime,
		event.Description,
		fmt.Sprintf("%v", event.Duration),
		fmt.Sprintf("%v", event.NotifyBefore),
		event.UserID,
	)

	if err != nil {
		s.Fail(err.Error())
	}

	rows, err := result.RowsAffected()

	if err != nil {
		s.Fail(err.Error())
	}

	if rows != 1 {
		s.Fail("expected to affect 1 row")
	}

}

func parseDuration(durationStr string) (time.Duration, error) {
	parts := strings.Split(durationStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	return duration, nil
}
