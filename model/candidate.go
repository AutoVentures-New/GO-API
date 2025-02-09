package model

import (
	"time"
)

type Candidate struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Status    Status    `json:"status"`
	Phone     string    `json:"phone"`
	BirthDate time.Time `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CandidateQuestionnaire struct {
	ID          int64                          `json:"id"`
	CandidateID int64                          `json:"candidate_id"`
	Type        string                         `json:"type"`
	Answers     []CandidateQuestionnaireAnswer `json:"answers"`
	ExpiredAt   time.Time                      `json:"expired_at"`
	CreatedAt   time.Time                      `json:"created_at"`
	UpdatedAt   time.Time                      `json:"updated_at"`
}

type CandidateQuestionnaireAnswer struct {
	QuestionID int64  `json:"question_id"`
	Answer     string `json:"answer"`
}
