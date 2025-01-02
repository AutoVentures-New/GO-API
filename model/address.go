package model

import (
	"time"
)

type Address struct {
	ID           int64     `json:"id"`
	Address      string    `json:"address"`
	Address2     string    `json:"address_2"`
	Neighborhood string    `json:"neighborhood"`
	ZipCode      string    `json:"zip_code"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
