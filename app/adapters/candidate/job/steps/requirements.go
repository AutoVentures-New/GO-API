package steps

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func SaveRequirements(
	ctx context.Context,
	application model.Application,
) (model.Application, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return application, err
	}

	jobRequirements, err := getJobRequirements(ctx, application.JobID, application.CompanyID)
	if err != nil {
		return application, err
	}

	match, reproved := matchValue(
		application,
		jobRequirements,
	)

	now := time.Now().UTC()
	jobApplicationRequirement := model.JobApplicationRequirement{
		ApplicationID: application.ID,
		Items:         application.JobApplicationRequirementItem,
		MatchValue:    match,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	itemsString, err := json.Marshal(jobApplicationRequirement.Items)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal application items")

		return application, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO job_application_requirements(application_id,items,match_value,created_at,updated_at) 
					VALUES(?,?,?,?,?)`,
		jobApplicationRequirement.ApplicationID,
		itemsString,
		jobApplicationRequirement.MatchValue,
		jobApplicationRequirement.CreatedAt,
		jobApplicationRequirement.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert job application requirements")

		_ = dbTransaction.Rollback()

		return application, err
	}

	if err := updateApplication(ctx, dbTransaction, &application, reproved); err != nil {
		_ = dbTransaction.Rollback()

		return application, err
	}

	if err := dbTransaction.Commit(); err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit transaction")

		return application, err
	}

	return application, nil
}

func matchValue(
	application model.Application,
	requirement model.JobRequirement,
) (int64, bool) {
	mapJobRequirement := make(map[string]model.JobRequirementItem)
	match := 0
	reproved := false

	for _, item := range requirement.Items {
		mapJobRequirement[fmt.Sprintf("%d", item.ID)] = item
	}

	for _, item := range application.JobApplicationRequirementItem {
		value, ok := mapJobRequirement[fmt.Sprintf("%d", item.ID)]
		if !ok {
			continue
		}

		if value.Required && !item.Checked {
			reproved = true

			continue
		}

		if !item.Checked {
			continue
		}

		match++
	}

	match = (match * 100) / len(requirement.Items)

	if !reproved {
		reproved = int64(match) < requirement.MinMatch
	}

	return int64(match), reproved
}

func getJobRequirements(
	ctx context.Context,
	jobID int64,
	companyID int64,
) (model.JobRequirement, error) {
	jobRequirement := model.JobRequirement{}

	var itemsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT items,min_match 
				FROM job_requirements WHERE company_id = ? AND job_id = ?`,
		companyID,
		jobID,
	).Scan(
		&itemsJSON,
		&jobRequirement.MinMatch,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get job requirement")

		return jobRequirement, err
	}

	if err = json.Unmarshal(itemsJSON, &jobRequirement.Items); err != nil {
		logrus.WithError(err).Error("Error to unmarshal job requirement")

		return jobRequirement, err
	}

	return jobRequirement, nil
}

func updateApplication(
	ctx context.Context,
	dbTransaction *sql.Tx,
	application *model.Application,
	reproved bool,
) error {
	application.UpdatedAt = time.Now().UTC()

	if reproved {
		application.Status = model.REPROVED
	} else {
		for index, value := range application.Steps {
			if value == model.REQUIREMENTS && len(application.Steps) > index+1 {
				application.CurrentStep = application.Steps[index+1]

				break
			}
		}
	}

	_, err := dbTransaction.ExecContext(
		ctx,
		`UPDATE job_applications set current_step = ?, status = ?, updated_at = ? WHERE id = ?`,
		application.CurrentStep,
		application.Status,
		application.UpdatedAt,
		application.ID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update job application")

		return err
	}

	return nil
}
