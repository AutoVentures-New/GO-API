package public

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hubjob/api/app/adapters/sendgrid"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/database"
	"github.com/sirupsen/logrus"
)

var ErrEmailNotFound = errors.New("email not found")

func ForgotPassword(
	ctx context.Context,
	email string,
	execType string,
) error {
	userId, err := getId(ctx, email, execType)
	if err != nil {
		return err
	}

	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return err
	}

	now := time.Now().UTC()
	token := uuid.New().String()

	_, err = dbTransaction.ExecContext(
		ctx,
		`INSERT INTO forgot_password(user_id,token,type,expired_at,created_at) VALUES (?,?,?,?,?) 
				ON DUPLICATE KEY UPDATE token = ?, expired_at = ?, created_at = ?`,
		userId,
		token,
		execType,
		now.Add(15*time.Minute),
		now,
		token,
		now.Add(15*time.Minute),
		now,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to insert code")

		return err
	}

	htmlText := fmt.Sprintf(
		"<html><a href='%s/trocar-senha/%s/%s'>Click aqui para trocar sua senha</a> </html>",
		config.Config.FrontendURL,
		token,
		strings.ToLower(execType),
	)

	err = sendgrid.SendEmail(
		ctx,
		"Esqueci minha senha",
		htmlText,
		"",
		email,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to send email forgot password")

		return err
	}
	fmt.Println("Email forgot password: ", email, token)

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return err
	}

	return nil
}

func getId(
	ctx context.Context,
	email string,
	execType string,
) (int64, error) {
	var id int64

	if execType == "COMPANY" {
		err := database.Database.QueryRowContext(
			ctx,
			`SELECT id FROM users WHERE email = ?`,
			email,
		).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrEmailNotFound
		}

		if err != nil {
			logrus.WithError(err).Error("Error to get user")

			return 0, err
		}

		return id, nil
	}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id FROM candidates WHERE email = ?`,
		email,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrEmailNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get candidate")

		return 0, err
	}

	return id, nil
}
