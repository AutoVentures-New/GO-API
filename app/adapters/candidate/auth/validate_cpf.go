package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

var CpfAlreadyExists = errors.New("cpf already exists")

func ValidateCpf(
	ctx context.Context,
	cpf string,
) error {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM candidates WHERE cpf = ?`,
		cpf,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate candidate cpf")

		return err
	}

	if count > 0 {
		return CpfAlreadyExists
	}

	return nil
}
