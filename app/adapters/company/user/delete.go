package user

import (
	"context"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

var ErrUserNotFound = errors.New("user not found")

func DeleteUser(
	ctx context.Context,
	id int64,
	companyID int64,
) error {
	result, err := database.Database.ExecContext(
		ctx,
		`DELETE FROM users WHERE id = ? AND company_id = ?`,
		id,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to delete user")

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Error("Error to get rows affected")

		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
