package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/hubjob/api/app/adapters/sendgrid"
	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

var EmailAlreadyExists = errors.New("Email already exists")

func SendEmailValidation(
	ctx context.Context,
	email string,
) error {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM users WHERE email = ?`,
		email,
	).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.WithError(err).Error("Error to validate company user email")

		return err
	}

	if count > 0 {
		return EmailAlreadyExists
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return err
	}

	code := generateRandomNumber()
	now := time.Now().UTC()

	_, err = database.Database.ExecContext(
		ctx,
		`INSERT INTO email_validations(email, code, created_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE code = ?, created_at = ?`,
		email,
		code,
		now,
		code,
		now,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to insert code")

		return err
	}

	htmlText := fmt.Sprintf(
		"<html>Code: %s <a href='https://google.com'>google</a> </html>",
		code,
	)

	err = sendgrid.SendEmail(
		ctx,
		"Codigo de validação de email",
		htmlText,
		"",
		email,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to send email validation")

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

func generateRandomNumber() string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	randomNumber := random.Intn(900000) + 100000

	return strconv.Itoa(randomNumber)
}
