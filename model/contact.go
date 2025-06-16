package model

type Contact struct {
	RecordNumber uint64  `json:"record_number"`
	Ulid         string  `json:"ulid"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	CompanyName  *string `json:"company_name"`
	Image        *string `json:"image"`
}
