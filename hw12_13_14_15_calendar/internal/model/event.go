package model

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrDateBusy      = errors.New("date is busy for this event")
	ErrEventNotFound = errors.New("event not found")
)

type Event struct {
	ID           uuid.UUID     `db:"id"`
	Title        string        `db:"title"`
	StartTime    time.Time     `db:"start_time"`
	Duration     time.Duration `db:"duration"`
	Description  string        `db:"description"`
	UserID       string        `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
	Sent         bool          `db:"sent"`
}

type Notification struct {
	EventID uuid.UUID
	Title   string
	Date    time.Time
	UserID  string
}
