package request

import "time"

type ContactDataQuery struct {
	ContactULID *string    `query:"contact_ulid"`
	Type        string     `query:"type"`
	Unread      *bool      `query:"unread"`
	From        *string    `query:"from"`
	To          *string    `query:"to"`
	Subject     *string    `query:"subject"`
	Folder      *string    `query:"folder"`
	DateFrom    *time.Time `query:"date_from"`
	DateTo      *time.Time `query:"date_to"`
	OrderBy     *string    `query:"order_by"`
	Sort        *string    `query:"sort"`
	Page        int        `query:"page"`
	Limit       int        `query:"limit"`
}
