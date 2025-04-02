package model

import "time"

type QueueJobStatus string

const (
	PENDING_JOB QueueJobStatus = "PENDING"
	PROCESSING  QueueJobStatus = "PROCESSING"
	ERROR       QueueJobStatus = "ERROR"
	FINISHED    QueueJobStatus = "FINISHED"
)

type QueueJob struct {
	ID             int64          `json:"id"`
	Type           string         `json:"type"`
	Status         QueueJobStatus `json:"status"`
	Configurations Configurations `json:"configurations"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"update_at"`
}

type Configurations struct {
	CandidateID int64 `json:"candidate_id,omitempty"`
	UserID      int64 `json:"user_id,omitempty"`
	CompanyID   int64 `json:"company_id,omitempty"`
}
