package profile

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateInfo(
	ctx context.Context,
	candidate model.Candidate,
) (model.Candidate, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return model.Candidate{}, err
	}

	candidate.UpdatedAt = time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE candidates set
    				name = ?, 
    				cpf = ?,  
    				phone = ?,  
    				birth_date = ?,  
    				updated_at = ?
				WHERE id = ?`,
		candidate.Name,
		candidate.CPF,
		candidate.Phone,
		candidate.BirthDate,
		candidate.UpdatedAt,
		candidate.ID,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to update candidate")

		return model.Candidate{}, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE candidate_addresses set
    				city = ?, 
    				state = ?,  
    				updated_at = ?
				WHERE candidate_id = ?`,
		candidate.Address.City,
		candidate.Address.State,
		candidate.UpdatedAt,
		candidate.ID,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to update candidate address")

		return model.Candidate{}, err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return model.Candidate{}, err
	}

	return candidate, nil
}
