package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrCandidateNotFound = errors.New("not found")

func GetCandidate(
	ctx context.Context,
	email string,
) (model.Candidate, error) {
	var candidate model.Candidate

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,name,cpf,email,password,status,phone,birth_date,created_at,updated_at FROM candidates WHERE email = ?`,
		email,
	).Scan(
		&candidate.ID,
		&candidate.Name,
		&candidate.CPF,
		&candidate.Email,
		&candidate.Password,
		&candidate.Status,
		&candidate.Phone,
		&candidate.BirthDate,
		&candidate.CreatedAt,
		&candidate.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return candidate, ErrCandidateNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to check candidate email")

		return candidate, err
	}

	_ = database.Database.QueryRowContext(
		ctx,
		`SELECT state, city FROM candidate_addresses WHERE candidate_id = ?`,
		candidate.ID,
	).Scan(
		&candidate.Address.State,
		&candidate.Address.City,
	)

	if candidate.Name == "" ||
		candidate.CPF == "" ||
		candidate.Email == "" ||
		candidate.Phone == "" ||
		candidate.Address.State == "" ||
		candidate.Address.City == "" {
		candidate.NeedCompleteProfile = true
	}

	if !candidate.NeedCompleteProfile {
		var bucketName string

		err := database.Database.QueryRowContext(
			ctx,
			`SELECT bucket_name,photo_path 
				FROM candidate_photos WHERE candidate_id = ?`,
			candidate.ID,
		).Scan(
			&bucketName,
		)
		if errors.Is(err, sql.ErrNoRows) {
			candidate.NeedCompleteProfile = true
		}
	}

	return candidate, nil
}
