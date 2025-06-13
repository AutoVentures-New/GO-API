package model

import (
	"encoding/json"
	"time"
)

type Note struct {
	RecordNumber   uint64          `json:"record_number"`
	Ulid           string          `json:"ulid"`
	CreatedBy      string          `json:"created_by"`
	To             EmailAddresses  `json:"to"`
	Subject        string          `json:"subject"`
	Done           bool            `json:"done"`
	Text           string          `json:"text"`
	CommentedAt    *string         `json:"commented_at"`
	Files          json.RawMessage `json:"files"`
	UserCreateDate time.Time       `json:"user_create_date"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
