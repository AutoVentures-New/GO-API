package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

func IsValidCode(
	ctx context.Context,
	code string,
	email string,
) (bool, error) {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM email_validations WHERE email = ? AND code = ?`,
		email,
		code,
	).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to select email and code")

		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
