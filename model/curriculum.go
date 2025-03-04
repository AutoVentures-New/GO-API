package model

import "time"

type Curriculum struct {
	ID               int64      `json:"id"`
	CandidateID      int64      `json:"candidate_id"`
	Gender           string     `json:"gender"`
	GenderIdentifier string     `json:"gender_identifier"`
	Color            string     `json:"color"`
	IsSpecialNeeds   bool       `json:"special_needs"`
	Languages        []Language `json:"languages"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	ProfessionalExperiences []ProfessionalExperience `json:"professional_experiences"`
	AcademicExperiences     []AcademicExperience     `json:"academic_experiences"`
}

type Language struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type ProfessionalExperience struct {
	ID          int64      `json:"id"`
	CandidateID int64      `json:"candidate_id"`
	Title       string     `json:"title"`
	Company     string     `json:"company"`
	AreaID      int64      `json:"area_id"`
	City        string     `json:"city"`
	State       string     `json:"state"`
	JobMode     string     `json:"job_mode"`
	CurrentJob  bool       `json:"current_job"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type AcademicExperience struct {
	ID          int64      `json:"id"`
	CandidateID int64      `json:"candidate_id"`
	Title       string     `json:"title"`
	Company     string     `json:"company"`
	AreaID      int64      `json:"area_id"`
	Status      string     `json:"status"`
	Level       string     `json:"level"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
