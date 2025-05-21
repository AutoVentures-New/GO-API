package public

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

var ErrUserTokenNotFound = errors.New("token not found")

func CreatePassword(
	ctx context.Context,
	token string,
	password string,
) (string, error) {
	hashPassword, err := pkg.HashPassword(password)
	if err != nil {
		logrus.WithError(err).Error("Error to hash password")

		return "", err
	}

	var emailValidation model.EmailValidation

	err = database.Database.QueryRowContext(
		ctx,
		`SELECT email,code FROM email_validations WHERE code = ?`,
		token,
	).Scan(
		&emailValidation.Email,
		&emailValidation.Code,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrUserTokenNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get token")

		return "", err
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return "", err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE users set password = ?,updated_at = ? WHERE email = ?`,
		hashPassword,
		time.Now().UTC(),
		emailValidation.Email,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to update users")

		return "", err
	}

	err = deleteEmailValidation(ctx, token, dbTransaction)
	if err != nil {
		_ = dbTransaction.Rollback()

		return "", err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return "", err
	}

	return emailValidation.Email, nil
}

func deleteEmailValidation(ctx context.Context, token string, dbTransaction *sql.Tx) error {
	_, err := dbTransaction.ExecContext(
		ctx,
		`DELETE FROM email_validations WHERE code = ?`,
		token,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to delete email validation")

		return err
	}

	return nil
}
