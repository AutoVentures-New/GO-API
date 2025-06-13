package model

import "time"

type ContactDataContact struct {
	Ulid            string    `json:"ulid"`
	ContactDataUlid string    `json:"contact_data_ulid"`
	ContactUlid     string    `json:"contact_ulid"`
	EmailUlid       string    `json:"email_ulid"`
	IsNew           bool      `json:"is_new"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
