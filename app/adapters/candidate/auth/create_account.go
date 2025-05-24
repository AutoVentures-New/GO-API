package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/hubjob/api/pkg"
	"github.com/sirupsen/logrus"
)

func CreateAccount(
	ctx context.Context,
	candidate model.Candidate,
	code string,
) (model.Candidate, error) {
	hashPassword, err := pkg.HashPassword(candidate.Password)
	if err != nil {
		logrus.WithError(err).Error("Error to hash password")

		return model.Candidate{}, err
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return model.Candidate{}, err
	}

	candidate.Status = model.ACTIVE
	candidate.Password = hashPassword
	candidate.CreatedAt = time.Now().UTC()
	candidate.UpdatedAt = candidate.CreatedAt

	result, err := dbTransaction.ExecContext(
		ctx,
		`INSERT INTO candidates(name,cpf,email,password,status,phone,birth_date,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`,
		candidate.Name,
		candidate.CPF,
		candidate.Email,
		candidate.Password,
		candidate.Status,
		candidate.Phone,
		candidate.BirthDate,
		candidate.CreatedAt,
		candidate.UpdatedAt,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to insert candidate")

		return model.Candidate{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to get last insert candidate id")

		return model.Candidate{}, err
	}

	candidate.Address.CreatedAt = candidate.CreatedAt
	candidate.Address.UpdatedAt = candidate.CreatedAt

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO candidate_addresses(candidate_id,city,state,created_at,updated_at) VALUES(?,?,?,?,?)`,
		lastInsertID,
		candidate.Address.City,
		candidate.Address.State,
		candidate.CreatedAt,
		candidate.UpdatedAt,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to insert candidate address")

		return model.Candidate{}, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO candidate_curriculum(candidate_id,gender,is_special_needs,languages,created_at,updated_at) 
				VALUES(?,?,?,?,?,?)`,
		lastInsertID,
		"",
		false,
		"[]",
		candidate.CreatedAt,
		candidate.UpdatedAt,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to insert candidate curriculum")

		return model.Candidate{}, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`DELETE FROM email_validations WHERE email = ? AND code = ?`,
		candidate.Email,
		code,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to delete email validation")

		return model.Candidate{}, err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return model.Candidate{}, err
	}

	candidate.ID = lastInsertID

	return candidate, nil
}

func CheckAlreadyExist(
	ctx context.Context,
	candidate model.Candidate,
) error {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM candidates WHERE email = ? OR cpf = ?`,
		candidate.Email,
		candidate.CPF,
	).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to check if candidate already exist")

		return err
	}

	if count != 0 {
		return ErrCandidateAlreadyExists
	}

	return nil
}
