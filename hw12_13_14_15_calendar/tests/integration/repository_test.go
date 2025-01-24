///go:build integration

package integration

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/model"
	sqlstorage "github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type Repository interface {
	CreateEvent(ctx context.Context, event model.Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, event model.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvents(ctx context.Context, date time.Time, offset int) ([]model.Event, error)
}

type IntegrationSuite struct {
	suite.Suite
	db *sqlx.DB
	r  Repository
}

func TestStorageIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) SetupSuite() {
	db, err := sqlx.ConnectContext(context.Background(), "postgres", "postgres://test:test@pg:5432/calendar?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
	s.r = sqlstorage.New(s.db)
}

func (s *IntegrationSuite) TearDownSuite() {
	_, _ = s.db.ExecContext(context.Background(), "TRUNCATE TABLE event")
	s.db.Close()
}

func (s *IntegrationSuite) SetupTest() {
	_, err := s.db.ExecContext(context.Background(), "TRUNCATE TABLE event")
	s.Require().NoError(err)
}

func (s *IntegrationSuite) TestCreateEvent() {
	const (
		eventTitle = "test create event"
	)

	startTime := time.Now()
	m := model.Event{
		Title:     eventTitle,
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "user1",
	}

	eventID, err := s.r.CreateEvent(context.Background(), m)
	s.Require().NoError(err)
	s.Require().NotEmpty(eventID)

	dbItem := s.getDirectItem(eventTitle)
	s.Require().NotEmpty(dbItem)
	s.Require().Equal(eventID, dbItem.ID)

	s.Require().Equal(m.StartTime.Format(time.DateTime), dbItem.StartTime.Format(time.DateTime))
	s.Require().Equal(m.Duration, dbItem.Duration)
	s.Require().Equal(m.UserID, dbItem.UserID)
}

func (s *IntegrationSuite) TestGetEvents() {
	const (
		eventTitle = "test get events"
	)

	startTime := time.Now()

	m := model.Event{
		Title:     eventTitle,
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "1000",
	}
	events, err := s.r.GetEvents(context.Background(), startTime, 0)
	s.Require().NoError(err)
	oldLen := len(events)

	dbID := s.createDirectItem(m)

	events, err = s.r.GetEvents(context.Background(), startTime, 0)
	createdEvent := events[len(events)-1]
	s.Require().NoError(err)
	s.Require().NotEmpty(events)
	s.Require().Equal(oldLen+1, len(events))
	s.Require().Equal(dbID, createdEvent.ID)
	s.Require().Equal(m.StartTime.Format(time.DateTime), createdEvent.StartTime.Format(time.DateTime))
	s.Require().Equal(m.Duration, createdEvent.Duration)
	s.Require().Equal(m.UserID, createdEvent.UserID)
}

func (s *IntegrationSuite) TestGetNotFound() {
	const (
		eventTitle = "test not found events"
	)

	startTime := time.Now()

	m := model.Event{
		Title:     eventTitle,
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "1000",
	}
	_ = s.createDirectItem(m)

	events, err := s.r.GetEvents(context.Background(), startTime.Add(time.Hour*60), 0)
	s.Require().NoError(err)
	s.Require().Empty(events)
}

func (s *IntegrationSuite) TestUpdateEvent() {
	const (
		eventTitleOld = "old event title"
		eventTitleNew = "new event title"
	)

	startTime := time.Now()
	m1 := model.Event{
		Title:       eventTitleOld,
		Description: "some description",
		StartTime:   startTime,
		Duration:    time.Hour,
		UserID:      "1000",
	}
	m2 := model.Event{
		Title:       eventTitleNew,
		Description: "updated description",
		StartTime:   startTime,
		Duration:    time.Hour,
	}

	dbID := s.createDirectItem(m1)
	err := s.r.UpdateEvent(context.Background(), dbID, m2)
	updatedEvent := s.getDirectItem(m2.Title)

	s.Require().NoError(err)
	s.Require().Equal(dbID, updatedEvent.ID)
	s.Require().Equal(m2.Title, updatedEvent.Title)
	s.Require().Equal(m2.Duration, updatedEvent.Duration)
	s.Require().Equal(m2.Description, updatedEvent.Description)
}

func (s *IntegrationSuite) TestDeleteEvent() {
	const (
		eventTitle = "test delete event 3"
	)
	startTime := time.Now()
	m := model.Event{
		Title:     eventTitle,
		StartTime: startTime,
		Duration:  time.Hour,
		UserID:    "1000",
	}
	dbID := s.createDirectItem(m)
	err := s.r.DeleteEvent(context.Background(), dbID)
	s.Require().NoError(err)

	events := s.getDirectItem(eventTitle)
	s.Require().Empty(events)
}

func (s *IntegrationSuite) createDirectItem(event model.Event) uuid.UUID {
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

	return eventId
}

func (s *IntegrationSuite) getDirectItem(title string) model.Event {
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
		return model.Event{}
	}

	duration, err := parseDuration(durationStr)
	if err != nil {
		s.Fail(err.Error())
	}
	e.Duration = duration

	return e
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
