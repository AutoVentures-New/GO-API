package model

import (
	"encoding/json"
	"time"
)

type Comment struct {
	RecordNumber   uint64          `json:"record_number"`
	Ulid           string          `json:"ulid"`
	CreatedBy      string          `json:"created_by"`
	Text           string          `json:"text"`
	CommentedAt    *string         `json:"commented_at"`
	Files          json.RawMessage `json:"files"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	CreatedByName  string          `json:"created_by_name"`
	CreatedByImage string          `json:"created_by_image"`
}
