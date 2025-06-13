package request

type ContactDataQuery struct {
	ContactULID string `query:"contact_ulid"`
	Page        int    `query:"page"`
	Limit       int    `query:"limit"`
}
