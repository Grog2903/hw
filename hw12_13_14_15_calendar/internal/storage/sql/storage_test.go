package sqlstorage

import (
	"context"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Grog2903/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestDB(ctx context.Context, t *testing.T) (*Storage, func()) {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	rootDir := filepath.Join(cwd, "../../..")
	configPath := filepath.Join(rootDir, "configs", "config.yaml")

	cfg, _ := config.LoadConfig(configPath)

	testStorage := New()
	err = testStorage.Connect(ctx, *cfg)
	if err != nil {
		t.Errorf("Failed to connect to test database: %v", err)
	}

	_, err = testStorage.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY NOT NULL,
			title TEXT NOT NULL,
			start_time TIMESTAMP NOT NULL,
			description TEXT
		)
	`)
	if err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	cleanup := func() {
		_, err := testStorage.db.ExecContext(ctx, `DROP TABLE IF EXISTS events`)
		if err != nil {
			panic("failed to cleanup test database: " + err.Error())
		}
		testStorage.Close(ctx)
	}

	return testStorage, cleanup
}

func TestStorage_CreateEvent(t *testing.T) {
	ctx := context.Background()
	testStorage, cleanup := setupTestDB(ctx, t)
	defer cleanup()

	event := storage.Event{
		Title:       "Test Event",
		StartTime:   time.Now(),
		Duration:    time.Hour,
		Description: "This is a test event",
	}

	eventID, err := testStorage.CreateEvent(ctx, event)
	require.NoError(t, err, "unexpected error during event creation")
	require.NotEmpty(t, eventID, "event ID should not be empty")

	var count int
	err = testStorage.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM events WHERE id = $1", eventID).Scan(&count)
	require.NoError(t, err, "unexpected error during event lookup")
	require.Equal(t, 1, count, "expected exactly one event in the database")
}

func TestStorage_UpdateEvent(t *testing.T) {
	ctx := context.Background()
	testStorage, cleanup := setupTestDB(ctx, t)
	defer cleanup()

	event := storage.Event{
		Title:       "Initial Event",
		StartTime:   time.Now(),
		Duration:    time.Hour,
		Description: "This is an initial event",
	}

	eventID, err := testStorage.CreateEvent(ctx, event)
	require.NoError(t, err, "unexpected error during event creation")

	updatedEvent := storage.Event{
		ID:          eventID,
		Title:       "Updated Event",
		StartTime:   event.StartTime.Add(2 * time.Hour),
		Description: "This is an updated event",
	}

	err = testStorage.UpdateEvent(ctx, eventID, updatedEvent)
	require.NoError(t, err, "unexpected error during event update")

	var title string
	err = testStorage.db.QueryRowContext(ctx, "SELECT title FROM events WHERE id = $1", eventID).Scan(&title)
	require.NoError(t, err, "unexpected error during event lookup")
	require.Equal(t, "Updated Event", title, "event title was not updated correctly")
}

func TestStorage_DeleteEvent(t *testing.T) {
	ctx := context.Background()
	testStorage, cleanup := setupTestDB(ctx, t)
	defer cleanup()

	event := storage.Event{
		Title:       "Test Event",
		StartTime:   time.Now(),
		Duration:    time.Hour,
		Description: "This is a test event",
	}

	eventID, err := testStorage.CreateEvent(ctx, event)
	require.NoError(t, err, "unexpected error during event creation")

	err = testStorage.DeleteEvent(ctx, eventID)
	require.NoError(t, err, "unexpected error during event deletion")

	var count int
	err = testStorage.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM events WHERE id = $1", eventID).Scan(&count)
	require.NoError(t, err, "unexpected error during event lookup")
	require.Equal(t, 0, count, "expected zero events in the database after deletion")
}

func TestStorage_GetEvents(t *testing.T) {
	ctx := context.Background()
	testStorage, cleanup := setupTestDB(ctx, t)
	defer cleanup()
	now := time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC)

	event1 := storage.Event{
		Title:       "Event 1",
		StartTime:   now,
		Duration:    time.Hour,
		Description: "First event",
	}
	event2 := storage.Event{
		Title:       "Event 2",
		StartTime:   now.AddDate(0, 0, 1),
		Duration:    2 * time.Hour,
		Description: "Second event",
	}

	_, err := testStorage.CreateEvent(ctx, event1)
	require.NoError(t, err, "unexpected error during event creation")
	_, err = testStorage.CreateEvent(ctx, event2)
	require.NoError(t, err, "unexpected error during event creation")

	events, err := testStorage.GetEvents(ctx, now, 2)
	require.NoError(t, err, "unexpected error during getting events")
	require.Equal(t, 2, len(events), "expected 2 events")
	require.Equal(t, "Event 1", events[0].Title, "unexpected event order")
	require.Equal(t, "Event 2", events[1].Title, "unexpected event order")
}