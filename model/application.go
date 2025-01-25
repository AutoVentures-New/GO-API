package model

import "time"

// steps
const (
	REQUIREMENTS    = "REQUIREMENTS"
	JOB_QUESTIONS   = "JOB_QUESTIONS"
	CULTURAL_FIT    = "CULTURAL_FIT"
	QUESTIONNAIRE   = "QUESTIONNAIRE"
	CANDIDATE_VIDEO = "CANDIDATE_VIDEO"
)

// status
const (
	FILLING            = "FILLING"
	WAITING_EVALUATION = "WAITING_EVALUATION"
	REPROVED           = "REPROVED"
	APPROVED           = "APPROVED"
	CANCELED           = "CANCELED"
)

type Application struct {
	ID          int64     `json:"id"`
	CompanyID   int64     `json:"company_id"`
	JobID       int64     `json:"job_id"`
	CandidateID int64     `json:"candidate_id"`
	Steps       []string  `json:"steps"`
	CurrentStep string    `json:"current_step"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	JobApplicationRequirementItem []JobApplicationRequirementItem `json:"job_requirement_items"`
}

type JobApplicationRequirement struct {
	ID            int64                           `json:"id"`
	ApplicationID int64                           `json:"application_id"`
	Items         []JobApplicationRequirementItem `json:"items"`
	MatchValue    int64                           `json:"match_value"`
	CreatedAt     time.Time                       `json:"created_at"`
	UpdatedAt     time.Time                       `json:"updated_at"`
}

type JobApplicationRequirementItem struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
}
