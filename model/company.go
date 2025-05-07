package model

import "time"

type Company struct {
	ID          int64     `json:"id"`
	CNPJ        string    `json:"cnpj"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Logo        *string   `json:"logo"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"update_at"`
}
