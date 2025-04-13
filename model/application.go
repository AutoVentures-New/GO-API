package model

import (
	"time"
)

// steps
const (
	REQUIREMENTS               = "REQUIREMENTS"
	JOB_QUESTIONS              = "JOB_QUESTIONS"
	CULTURAL_FIT               = "CULTURAL_FIT"
	QUESTIONNAIRE_BEHAVIORAL   = "QUESTIONNAIRE_BEHAVIORAL"
	QUESTIONNAIRE_PROFESSIONAL = "QUESTIONNAIRE_PROFESSIONAL"
	CANDIDATE_VIDEO            = "CANDIDATE_VIDEO"
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
	Questions                     []ApplicationQuestion           `json:"questions"`
	CulturalFit                   *JobApplicationCulturalFit      `json:"cultural_fit"`
	JobVideoQuestions             *JobVideoQuestions              `json:"job_video_questions"`

	Candidate                    *Candidate                    `json:"candidate"`
	JobApplicationRequirement    *JobApplicationRequirement    `json:"job_application_requirement"`
	JobApplicationCandidateVideo *JobApplicationCandidateVideo `json:"job_application_candidate_video"`
	JobApplicationQuestion       *JobApplicationQuestion       `json:"job_application_question"`
	Job                          *Job                          `json:"job"`
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

type ApplicationQuestion struct {
	ID      int64               `json:"id"`
	Title   string              `json:"title"`
	Type    string              `json:"type"`
	Answers []ApplicationAnswer `json:"answers"`
}

type ApplicationAnswer struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Checked bool   `json:"checked"`
	Answer  string `json:"answer"`
}

type JobApplicationQuestion struct {
	ID            int64                 `json:"id"`
	ApplicationID int64                 `json:"application_id"`
	Questions     []ApplicationQuestion `json:"questions"`
	Score         int64                 `json:"score"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`

	JobQuestions []Question `json:"job_questions"`
}

type JobApplicationCulturalFit struct {
	ID            int64                             `json:"id"`
	ApplicationID int64                             `json:"application_id"`
	Answers       []JobApplicationCulturalFitAnswer `json:"answers"`
	MatchValue    int64                             `json:"match_value"`
	CreatedAt     time.Time                         `json:"created_at"`
	UpdatedAt     time.Time                         `json:"updated_at"`

	JobCulturalFit *JobCulturalFit `json:"job_cultural_fit"`
}

type JobApplicationCulturalFitAnswer struct {
	CulturalFitID int64  `json:"cultural_fit_id"`
	Answer        string `json:"answer"`
}

type JobApplicationCandidateVideo struct {
	ID            int64     `json:"id"`
	ApplicationID int64     `json:"application_id"`
	BucketName    string    `json:"bucket_name"`
	VideoPath     string    `json:"video_path"`
	Score         int64     `json:"score"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
