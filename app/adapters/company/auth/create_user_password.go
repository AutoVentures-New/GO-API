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

func CreateUserPassword(
	ctx context.Context,
	user model.User,
	code string,
) (model.User, error) {
	var email string

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT email FROM email_validations WHERE code = ?`,
		code,
	).Scan(&email)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to check if user code exists")

		return model.User{}, err
	}

	hashPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		logrus.WithError(err).Error("Error to hash password")

		return model.User{}, err
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return model.User{}, err
	}

	user.Email = email
	user.Status = model.ACTIVE
	user.Password = hashPassword
	user.UpdatedAt = time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE users SET password = ?, status = ?, updated_at = ? WHERE id = ? AND status = 'PENDING'`,
		user.Password,
		user.Status,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to update password user")

		return model.User{}, err
	}

	_, err = dbTransaction.ExecContext(
		ctx,
		`DELETE FROM email_validations WHERE email = ? AND code = ?`,
		user.Email,
		code,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to delete email validation")

		return model.User{}, err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return model.User{}, err
	}

	return user, nil
}
