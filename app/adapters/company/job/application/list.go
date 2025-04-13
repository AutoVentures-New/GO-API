package job

import (
	"context"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

const (
	CREATED_AT_DESC = "CREATED_AT_DESC"
	CREATED_AT_ASC  = "CREATED_AT_ASC"
	HIGHEST_SCORE   = "HIGHEST_SCORE"
	LOWEST_SCORE    = "LOWEST_SCORE"
)

type ListJobApplicationsRequest struct {
	FilterCandidateName string `json:"candidate_name"`
	OrderBy             string `json:"order_by"`
}

func ListJobApplications(
	ctx context.Context,
	companyID int64,
	jobID int64,
	request ListJobApplicationsRequest,
) ([]model.Application, error) {
	applications := make([]model.Application, 0)

	args := make([]any, 0)

	query := `SELECT 
    					j.id,
    					j.company_id,
    					j.job_id,
    					j.candidate_id,
    					j.current_step,
    					j.status,
    					j.created_at,
    					j.updated_at,
    					c.name,
    					r.match_value as r_match_value,
    					f.match_value as f_match_value,
    					v.score,
    					q.score
				FROM job_applications j 
				JOIN candidates c on c.id = j.candidate_id
				LEFT JOIN job_application_requirements r on j.id = r.application_id
				LEFT JOIN job_application_cultural_fit f on j.id = f.application_id
				LEFT JOIN job_application_candidate_videos v on j.id = v.application_id
				LEFT JOIN job_application_questions q on j.id = q.application_id
				WHERE j.company_id = ? AND j.job_id = ?`

	args = append(args, companyID)
	args = append(args, jobID)

	if len(request.FilterCandidateName) > 0 {
		query = query + " AND c.name LIKE ?"
		args = append(args, "%"+request.FilterCandidateName+"%")
	}

	query += " ORDER BY j.id DESC"

	rows, err := database.Database.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list job application")

		return applications, err
	}

	defer rows.Close()

	for rows.Next() {
		var application model.Application

		var candidate model.Candidate

		var requirementMatchValue, culturalFitMatchValue, candidateVideoScore, jobQuestionsScore *int64

		err = rows.Scan(
			&application.ID,
			&application.CompanyID,
			&application.JobID,
			&application.CandidateID,
			&application.CurrentStep,
			&application.Status,
			&application.CreatedAt,
			&application.UpdatedAt,

			&candidate.Name,
			&requirementMatchValue,
			&culturalFitMatchValue,
			&candidateVideoScore,
			&jobQuestionsScore,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to unmarshal application")

			return applications, err
		}

		application.Candidate = &candidate
		application.JobApplicationRequirement = &model.JobApplicationRequirement{}
		application.CulturalFit = &model.JobApplicationCulturalFit{}
		application.JobApplicationCandidateVideo = &model.JobApplicationCandidateVideo{}
		application.JobApplicationQuestion = &model.JobApplicationQuestion{}

		if requirementMatchValue != nil {
			application.JobApplicationRequirement.MatchValue = *requirementMatchValue
		}

		if culturalFitMatchValue != nil {
			application.CulturalFit.MatchValue = *culturalFitMatchValue
		}

		if candidateVideoScore != nil {
			application.JobApplicationCandidateVideo.Score = *candidateVideoScore
		}

		if jobQuestionsScore != nil {
			application.JobApplicationQuestion.Score = *jobQuestionsScore
		}

		applications = append(applications, application)
	}

	return applications, nil
}
