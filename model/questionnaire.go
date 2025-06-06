package model

import "time"

const SINGLE_CHOICE = "SINGLE_CHOICE"
const MULTIPLE_CHOICE = "MULTIPLE_CHOICE"
const OPEN_FIELD = "OPEN_FIELD"

type Questionnaire struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CompanyID int64      `json:"company_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	QuestionnaireID int64     `json:"questionnaire_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Answers         []Answer  `json:"answers"`
}

type Answer struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	QuestionnaireID int64     `json:"questionnaire_id"`
	QuestionID      int64     `json:"question_id"`
	IsCorrect       bool      `json:"is_correct"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ExampleQuestions struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Subs        []ExampleQuestionsSub `json:"subs"`
}

type ExampleQuestionsSub struct {
	Title     string   `json:"title"`
	Questions []string `json:"questions"`
}
