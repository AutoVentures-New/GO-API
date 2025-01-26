package get_application

import (
	"context"
	"encoding/json"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func GetJobRequirements(
	ctx context.Context,
	jobID int64,
	companyID int64,
) ([]model.JobApplicationRequirementItem, error) {
	appRequirements := make([]model.JobApplicationRequirementItem, 0)
	jobRequirement := model.JobRequirement{}

	var itemsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT items 
				FROM job_requirements WHERE company_id = ? AND job_id = ?`,
		companyID,
		jobID,
	).Scan(
		&itemsJSON,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get job requirement")

		return nil, err
	}

	if err = json.Unmarshal(itemsJSON, &jobRequirement.Items); err != nil {
		logrus.WithError(err).Error("Error to unmarshal job requirement")

		return nil, err
	}

	for _, jobRequirementItem := range jobRequirement.Items {
		appRequirements = append(appRequirements, model.JobApplicationRequirementItem{
			ID:      jobRequirementItem.ID,
			Name:    jobRequirementItem.Name,
			Checked: false,
		})
	}

	return appRequirements, nil
}
