package model

import (
	"time"
)

type CalendarEventFile struct {
	RecordNumber uint64    `json:"record_number"`
	Ulid         string    `json:"ulid"`
	EventUlid    *string   `json:"event_ulid"`
	IsExternal   bool      `json:"is_external"`
	Name         string    `json:"name"`
	Extension    string    `json:"extension"`
	Link         string    `json:"link"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
