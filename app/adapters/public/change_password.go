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

var ErrTokenNotFound = errors.New("token not found")

func ChangePassword(
	ctx context.Context,
	token string,
	password string,
) error {
	hashPassword, err := pkg.HashPassword(password)
	if err != nil {
		logrus.WithError(err).Error("Error to hash password")

		return err
	}

	var forgotPassword model.ForgotPassword

	err = database.Database.QueryRowContext(
		ctx,
		`SELECT id,user_id,token,type,expired_at,created_at FROM forgot_password WHERE token = ?`,
		token,
	).Scan(
		&forgotPassword.ID,
		&forgotPassword.UserID,
		&forgotPassword.Token,
		&forgotPassword.Type,
		&forgotPassword.ExpiredAt,
		&forgotPassword.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrTokenNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get token")

		return err
	}

	if forgotPassword.ExpiredAt.Before(time.Now().UTC()) {
		logrus.Error("Expired token")

		return ErrTokenNotFound
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return err
	}

	if forgotPassword.Type == "CANDIDATE" {
		_, err = dbTransaction.ExecContext(
			ctx,
			`UPDATE candidates set password = ?,updated_at = ? WHERE id = ?`,
			hashPassword,
			time.Now().UTC(),
			forgotPassword.UserID,
		)
		if err != nil {
			_ = dbTransaction.Rollback()

			logrus.WithError(err).Error("Error to update candidate")

			return err
		}
	} else {
		_, err = dbTransaction.ExecContext(
			ctx,
			`UPDATE users set password = ?,updated_at = ? WHERE id = ?`,
			hashPassword,
			time.Now().UTC(),
			forgotPassword.UserID,
		)
		if err != nil {
			_ = dbTransaction.Rollback()

			logrus.WithError(err).Error("Error to update users")

			return err
		}
	}

	err = deleteForgotPassword(ctx, forgotPassword.ID, dbTransaction)
	if err != nil {
		_ = dbTransaction.Rollback()

		return err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return err
	}

	return nil
}

func deleteForgotPassword(ctx context.Context, id int64, dbTransaction *sql.Tx) error {
	_, err := dbTransaction.ExecContext(
		ctx,
		`DELETE FROM forgot_password WHERE id = ?`,
		id,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to delete forgot password")

		return err
	}

	return nil
}
