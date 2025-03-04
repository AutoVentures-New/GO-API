package job

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListJobApplications(
	ctx context.Context,
	companyID int64,
	jobID int64,
) ([]model.Application, error) {
	applications := make([]model.Application, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT 
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
    					v.score
				FROM job_applications j 
				JOIN candidates c on c.id = j.candidate_id
				LEFT JOIN job_application_requirements r on j.id = r.application_id
				LEFT JOIN job_application_cultural_fit f on j.id = f.application_id
				LEFT JOIN job_application_candidate_videos v on j.id = v.application_id
				WHERE j.company_id = ? AND j.job_id = ?
				ORDER BY id DESC`,
		companyID,
		jobID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list job application")

		return applications, err
	}

	defer rows.Close()

	for rows.Next() {
		var application model.Application

		var candidate model.Candidate

		var requirementMatchValue, culturalFitMatchValue, candidateVideoScore *int64

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
		)
		if err != nil {
			logrus.WithError(err).Error("Error to unmarshal application")

			return applications, err
		}

		application.Candidate = &candidate
		application.JobApplicationRequirement = &model.JobApplicationRequirement{}
		application.CulturalFit = &model.JobApplicationCulturalFit{}
		application.JobApplicationCandidateVideo = &model.JobApplicationCandidateVideo{}

		if requirementMatchValue != nil {
			application.JobApplicationRequirement.MatchValue = *requirementMatchValue
		}

		if culturalFitMatchValue != nil {
			application.CulturalFit.MatchValue = *culturalFitMatchValue
		}

		if candidateVideoScore != nil {
			application.JobApplicationCandidateVideo.Score = *candidateVideoScore
		}

		applications = append(applications, application)
	}

	return applications, nil
}
