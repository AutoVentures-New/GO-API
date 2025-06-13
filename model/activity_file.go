package model

import (
	"encoding/json"
	"time"
)

type ActivityFile struct {
	RecordNumber   uint64          `json:"record_number"`
	Ulid           string          `json:"ulid"`
	CreatedBy      string          `json:"created_by"`
	To             EmailAddresses  `json:"to"`
	Subject        string          `json:"subject"`
	Done           bool            `json:"done"`
	Files          json.RawMessage `json:"files"`
	UserCreateDate time.Time       `json:"user_create_date"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Comments       []Comment       `json:"comments"`
}

func (a ActivityFile) GetUlid() string {
	return a.Ulid

}
