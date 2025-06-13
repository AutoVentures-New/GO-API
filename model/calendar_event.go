package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type CalendarEvent struct {
	RecordNumber      uint64              `json:"record_number"`
	Ulid              string              `json:"ulid"`
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	Participants      EmailAddresses      `json:"participants"`
	When              json.RawMessage     `json:"when"`
	Location          *string             `json:"location"`
	Recurrence        json.RawMessage     `json:"recurrence"`
	Notifications     json.RawMessage     `json:"notifications"`
	Conferencing      json.RawMessage     `json:"conferencing"`
	ConferenceRecords sql.NullString      `json:"conference_records"`
	OrganizerName     string              `json:"organizer_name"`
	OrganizerEmail    string              `json:"organizer_email"`
	Owner             string              `json:"owner"`
	Done              bool                `json:"done"`
	StartDate         *time.Time          `json:"start_date"`
	EndDate           *time.Time          `json:"end_date"`
	AllDay            bool                `json:"all_day"`
	Type              string              `json:"type"`
	Sequence          uint                `json:"sequence"`
	Files             json.RawMessage     `json:"files"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	ActivityFiles     []CalendarEventFile `json:"activity_files"`
}

func (c CalendarEvent) GetUlid() string {
	return c.Ulid
}
