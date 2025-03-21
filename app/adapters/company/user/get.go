package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func GetUser(
	ctx context.Context,
	id int64,
	companyID int64,
) (model.User, error) {
	var user model.User

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,name,cpf,email,status,company_id,created_at,updated_at FROM users WHERE id = ? AND company_id = ?`,
		id,
		companyID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.CPF,
		&user.Email,
		&user.Status,
		&user.CompanyID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrUserNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get user")

		return user, err
	}

	return user, nil
}
