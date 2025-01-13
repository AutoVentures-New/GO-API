package model

import "time"

type Benefit struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CompanyID int64     `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
