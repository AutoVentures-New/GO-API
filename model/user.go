package model

import (
	"time"
)

type Role string

const (
	ADMINISTRATOR Role = "ADMINISTRATOR"
	MANAGER       Role = "MANAGER"
	RECRUITER     Role = "RECRUITER"
)

type User struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Phone     *string   `json:"phone"`
	Role      Role      `json:"role"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Company *Company `json:"company"`
}
