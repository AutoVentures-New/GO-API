package model

type User struct {
	Account   string  `json:"account"`
	User      string  `json:"user"`
	Ulid      string  `json:"ulid"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Image     *string `json:"image"`
}
