package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type EmailBucket struct {
	RecordNumber     uint64          `json:"record_number"`
	Ulid             string          `json:"ulid"`
	AccountID        sql.NullString  `json:"account_id"`
	MessageID        string          `json:"message_id"`
	ThreadID         string          `json:"thread_id"`
	Subject          string          `json:"subject"`
	From             EmailAddresses  `json:"from"`
	To               EmailAddresses  `json:"to"`
	CC               EmailAddresses  `json:"cc"`
	BCC              EmailAddresses  `json:"bcc"`
	ReplyTo          json.RawMessage `json:"reply_to"`
	Headers          json.RawMessage `json:"headers"`
	Starred          bool            `json:"starred"`
	Unread           bool            `json:"unread"`
	ReplyToMessageID sql.NullString  `json:"reply_to_message_id"`
	Body             string          `json:"body"`
	Files            json.RawMessage `json:"files"`
	Folder           string          `json:"folder"`
	Links            uint            `json:"links"`
	Opens            uint            `json:"opens"`
	LinkClicks       json.RawMessage `json:"link_clicks"`
	IsTracked        bool            `json:"is_tracked"`
	Date             time.Time       `json:"date"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}
