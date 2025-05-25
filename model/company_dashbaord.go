package model

import "time"

type OpenJobs struct {
	Count         int64 `json:"count"`
	CloseToFinish int64 `json:"close_to_finish"`
	Jobs          []struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	} `json:"jobs"`
}

type Applications struct {
	Count                int64             `json:"count"`
	WaitingEvaluation    int64             `json:"waiting_evaluation"`
	Applications         []ApplicationDash `json:"applications"`
	ApplicationDashDates map[string]int64  `json:"application_dash_dates"`
}

type ApplicationDash struct {
	ID            int64     `json:"id"`
	JobID         int64     `json:"job_id"`
	Status        string    `json:"status"`
	CandidateID   string    `json:"candidate_id"`
	CandidateName string    `json:"candidate_name"`
	CreatedAt     time.Time `json:"created_at"`
}
