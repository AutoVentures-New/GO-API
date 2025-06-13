package model

import (
	"time"
)

type Call struct {
	RecordNumber       uint64         `json:"record_number"`
	Ulid               string         `json:"ulid"`
	Subject            string         `json:"subject"`
	CreatedBy          string         `json:"created_by"`
	UserPhoneNumber    string         `json:"user_phone_number"`
	ContactPhoneNumber string         `json:"contact_phone_number"`
	Done               bool           `json:"done"`
	To                 EmailAddresses `json:"to"`
	CallType           string         `json:"call_type"`
	Direction          string         `json:"direction"`
	Outcome            string         `json:"outcome"`
	Notes              string         `json:"notes"`
	UserCreateDate     time.Time      `json:"user_create_date"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}
