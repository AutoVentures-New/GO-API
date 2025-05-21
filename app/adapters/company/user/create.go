package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hubjob/api/app/adapters/sendgrid"
	"github.com/hubjob/api/config"
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

	user.Password = model.NULL_PASSWORD
	user.Status = model.PENDING
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return user, err
	}

	result, err := dbTransaction.ExecContext(
		ctx,
		`INSERT INTO users(name,cpf,phone,email,password,status,company_id,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`,
		user.Name,
		user.CPF,
		user.Phone,
		user.Email,
		user.Password,
		user.Status,
		user.CompanyID,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to insert user")

		return user, err
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to get last insert user id")

		return user, err
	}

	err = sendEmailValidation(ctx, user.Email, dbTransaction)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to send email validation")

		return user, err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return user, err
	}

	return user, nil
}

func sendEmailValidation(
	ctx context.Context,
	email string,
	dbTransaction *sql.Tx,
) error {
	token := uuid.New().String()
	now := time.Now().UTC()

	_, err := dbTransaction.ExecContext(
		ctx,
		`INSERT INTO email_validations(email, code, created_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE code = ?, created_at = ?`,
		email,
		token,
		now,
		token,
		now,
	)
	if err != nil {
		return err
	}

	htmlText := fmt.Sprintf(
		"<html><a href='%s/criar-conta/%s'>Click aqui para criar sua senha</a> </html>",
		config.Config.FrontendURL,
		token,
	)

	err = sendgrid.SendEmail(
		ctx,
		"Crie sua senha",
		htmlText,
		"",
		email,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to send email validation")

		return err
	}
	fmt.Println("Email company user validation code: ", email, token)

	return nil
}
