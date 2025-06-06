package model

import (
	"time"
)

const (
	BEHAVIORAL   = "BEHAVIORAL"
	PROFESSIONAL = "PROFESSIONAL"
	NONE         = "NONE"
)

type Job struct {
	ID                  int64      `json:"id"`
	Title               string     `json:"title"`
	CompanyID           int64      `json:"company_id"`
	AreaID              int64      `json:"area_id"`
	IsTalentBank        bool       `json:"is_talent_bank"`
	IsSpecialNeeds      bool       `json:"is_special_needs"`
	Description         string     `json:"description"`
	JobMode             string     `json:"job_mode"`
	ContractingModality string     `json:"contracting_modality"`
	State               string     `json:"state"`
	City                string     `json:"city"`
	Responsibilities    string     `json:"responsibilities"`
	Questionnaire       string     `json:"questionnaire"`
	VideoLink           string     `json:"video_link"`
	Status              string     `json:"status"`
	PublishAt           time.Time  `json:"publish_at"`
	FinishAt            *time.Time `json:"finish_at"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	JobCulturalFit *JobCulturalFit    `json:"cultural_fit"`
	JobRequirement *JobRequirement    `json:"requirements"`
	Benefits       []Benefit          `json:"benefits"`
	VideoQuestions *JobVideoQuestions `json:"video_questions"`
	Questions      []Question         `json:"questions"`
	Area           *Area              `json:"area"`

	CandidateHasApplication bool `json:"candidate_has_application"`
}

type JobCulturalFit struct {
	ID        int64                  `json:"id"`
	CompanyID int64                  `json:"company_id"`
	JobID     int64                  `json:"job_id"`
	Answers   []JobCulturalFitAnswer `json:"answers"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type JobCulturalFitAnswer struct {
	CulturalFitID int64  `json:"cultural_fit_id"`
	Answer        string `json:"answer"`
}

type JobRequirement struct {
	ID        int64                `json:"id"`
	CompanyID int64                `json:"company_id"`
	JobID     int64                `json:"job_id"`
	Items     []JobRequirementItem `json:"items"`
	MinMatch  int64                `json:"min_match"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type JobRequirementItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type JobVideoQuestions struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"`
	JobID     int64     `json:"job_id"`
	Questions []string  `json:"questions"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Area struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type State struct {
	Name   string `json:"name"`
	Cities []City `json:"cities"`
}

type City struct {
	Name string `json:"name"`
}

type JobVideo struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"`
	VideoLink string    `json:"video_link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
