package model

import "time"

type ContactData struct {
	Ulid         string         `json:"ulid"`
	ProviderUlid string         `json:"provider_ulid"`
	Type         string         `json:"type"`
	Identifier   string         `json:"identifier"`
	SearchField  string         `json:"search_field"`
	From         EmailAddresses `json:"from"`
	To           EmailAddresses `json:"to"`
	CC           EmailAddresses `json:"cc"`
	BCC          EmailAddresses `json:"bcc"`
	Folder       string         `json:"folder"`
	Unread       bool           `json:"unread"`
	Date         time.Time      `json:"date"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
