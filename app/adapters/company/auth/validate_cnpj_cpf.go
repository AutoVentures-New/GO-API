package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

var CnpjAlreadyExists = errors.New("cnpj already exists")
var CpfAlreadyExists = errors.New("cpf already exists")

func ValidateCnpjCpf(
	ctx context.Context,
	cnpj string,
	cpf string,
) error {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM companies WHERE cnpj = ?`,
		cnpj,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate company cnpj")

		return err
	}

	if count > 0 {
		return CnpjAlreadyExists
	}

	err = database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM users WHERE cpf = ?`,
		cpf,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate company user cpf")

		return err
	}

	if count > 0 {
		return CpfAlreadyExists
	}

	return nil
}
