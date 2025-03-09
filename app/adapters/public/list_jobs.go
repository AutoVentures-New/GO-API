package public

import (
	"context"
	"strings"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

type Filter struct {
	CompanyID int    `query:"company_id"`
	AreaID    int    `query:"area_id"`
	State     string `query:"state"`
	City      string `query:"city"`
	Page      int64  `query:"page"`
	Size      int64  `query:"size"`
}

func ListJobs(
	ctx context.Context,
	filter Filter,
) ([]model.Job, int64, error) {
	jobs := make([]model.Job, 0)

	filter.Page--

	queryCount := `SELECT count(*) FROM jobs`

	query := `SELECT 
    			id,
    			title,
    			company_id,
    			area_id,
    			is_talent_bank,
    			is_special_needs,
    			description,
    			job_mode,
    			contracting_modality,
    			state,
    			city,
    			responsibilities,
    			questionnaire,
    			video_link,
    			status,
    			publish_at,
    			finish_at,
    			created_at,
    			updated_at 
			FROM jobs`

	where := make([]string, 0)
	attributes := make([]any, 0)

	where = append(where, `status = ?`)
	attributes = append(attributes, model.ACTIVE)

	if filter.CompanyID > 0 {
		where = append(where, `company_id = ?`)
		attributes = append(attributes, filter.CompanyID)
	}

	if filter.AreaID > 0 {
		where = append(where, `area_id = ?`)
		attributes = append(attributes, filter.AreaID)
	}

	if len(filter.State) > 0 {
		where = append(where, `state = ?`)
		attributes = append(attributes, filter.State)
	}

	if len(filter.City) > 0 {
		where = append(where, `city = ?`)
		attributes = append(attributes, filter.City)
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
		queryCount += " WHERE " + strings.Join(where, " AND ")
	}

	var count int64

	err := database.Database.QueryRowContext(
		ctx,
		queryCount,
		attributes...,
	).Scan(
		&count,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get jobs count")

		return nil, 0, err
	}

	query += " ORDER BY title ASC LIMIT ? OFFSET ?"
	attributes = append(attributes, filter.Size)
	attributes = append(attributes, filter.Page*filter.Size)

	rows, err := database.Database.QueryContext(
		ctx,
		query,
		attributes...,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list jobs")

		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		job := model.Job{}
		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.CompanyID,
			&job.AreaID,
			&job.IsTalentBank,
			&job.IsSpecialNeeds,
			&job.Description,
			&job.JobMode,
			&job.ContractingModality,
			&job.State,
			&job.City,
			&job.Responsibilities,
			&job.Questionnaire,
			&job.VideoLink,
			&job.Status,
			&job.PublishAt,
			&job.FinishAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan jobs")

			return nil, 0, err
		}

		jobs = append(jobs, job)
	}

	return jobs, count, nil
}
