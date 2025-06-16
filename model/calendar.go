package model

import "time"

type Calendar struct {
	RecordNumber uint64    `json:"record_number"`
	Ulid         string    `json:"ulid"`
	AccountUlid  string    `json:"account_ulid"`
	UserUlid     string    `json:"user_ulid"`
	ProviderUlid *string   `json:"provider_ulid,omitempty"`
	ExternalID   *string   `json:"external_id,omitempty"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	Color        string    `json:"color"`
	Sequence     uint      `json:"sequence"`
	Checked      bool      `json:"checked"`
	Hide         bool      `json:"hide"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
