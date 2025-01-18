package model

import (
	"time"
)

type Job struct {
	ID                  int64     `json:"id"`
	Title               string    `json:"title"`
	CompanyID           int64     `json:"company_id"`
	IsTalentBank        bool      `json:"is_talent_bank"`
	IsSpecialNeeds      bool      `json:"is_special_needs"`
	Description         string    `json:"description"`
	JobMode             string    `json:"job_mode"`
	ContractingModality string    `json:"contracting_modality"`
	State               string    `json:"state"`
	City                string    `json:"city"`
	Responsibilities    string    `json:"responsibilities"`
	Questionnaire       string    `json:"questionnaire"`
	VideoLink           string    `json:"video_link"`
	Status              string    `json:"status"`
	PublishAt           time.Time `json:"publish_at"`
	FinishAt            time.Time `json:"finish_at"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	JobCulturalFit JobCulturalFit `json:"cultural_fit"`
	JobRequirement JobRequirement `json:"requirements"`
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
