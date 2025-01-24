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
	ID           uuid.UUID
	Title        string
	StartTime    time.Time
	Duration     time.Duration
	Description  string
	UserID       string
	NotifyBefore time.Duration
	Sent         bool
}

type Notification struct {
	EventID uuid.UUID
	Title   string
	Date    time.Time
	UserID  string
}
