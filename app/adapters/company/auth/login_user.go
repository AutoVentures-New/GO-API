package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func LoginUser(
	ctx context.Context,
	email string,
) (model.User, error) {
	var user model.User

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,name,cpf,email,password,status,company_id,created_at,updated_at FROM users WHERE email = ?`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.CPF,
		&user.Email,
		&user.Password,
		&user.Status,
		&user.CompanyID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return user, nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to check email and password")

		return user, err
	}

	return user, nil
}
