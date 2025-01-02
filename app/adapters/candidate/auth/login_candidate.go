package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func LoginCandidate(
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
		return candidate, nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to check candidate email")

		return candidate, err
	}

	return candidate, nil
}
