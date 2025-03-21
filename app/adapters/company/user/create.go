package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrUserEmailAlreadyExists = errors.New("user email already exists")

var ErrUserCPFAlreadyExists = errors.New("user cpf already exists")

func CreateUser(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM users WHERE email = ?`,
		user.Email,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate user email")

		return user, err
	}

	if count > 0 {
		return user, ErrUserEmailAlreadyExists
	}

	err = database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM users WHERE cpf = ?`,
		user.CPF,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate user cpf")

		return user, err
	}

	if count > 0 {
		return user, ErrUserCPFAlreadyExists
	}

	user.Password = "NULL_PASSWORD"
	user.Status = model.INACTIVE
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt

	result, err := database.Database.ExecContext(
		ctx,
		`INSERT INTO users(name,cpf,email,password,status,company_id,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)`,
		user.Name,
		user.CPF,
		user.Email,
		user.Password,
		user.Status,
		user.CompanyID,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert user")

		return user, err
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert user id")

		return user, err
	}

	return user, nil
}
