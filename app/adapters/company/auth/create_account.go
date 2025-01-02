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

func CreateAccount(
	ctx context.Context,
	user model.User,
	company model.Company,
	code string,
) (model.User, error) {
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

	company, err = createCompany(ctx, dbTransaction, company)
	if err != nil {
		_ = dbTransaction.Rollback()

		return model.User{}, err
	}

	user.CompanyID = company.ID
	user.Status = model.ACTIVE
	user.Password = hashPassword
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt

	result, err := dbTransaction.ExecContext(
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
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to insert user")

		return model.User{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to get last insert user id")

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

	user.ID = lastInsertID

	return user, nil
}

func CheckAlreadyExist(
	ctx context.Context,
	user model.User,
	company model.Company,
) error {
	var count int

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM users WHERE email = ? OR cpf = ?`,
		user.Email,
		user.CPF,
	).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to check if user already exist")

		return err
	}

	if count != 0 {
		return ErrUserAlreadyExists
	}

	err = database.Database.QueryRowContext(
		ctx,
		`SELECT COUNT(0) FROM company WHERE cnpj = ?`,
		company.CNPJ,
	).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to check if company already exist")

		return err
	}

	if count != 0 {
		return ErrCompanyAlreadyExists
	}

	return nil
}

func createCompany(
	ctx context.Context,
	dbTransaction *sql.Tx,
	company model.Company,
) (model.Company, error) {
	company.Status = model.ACTIVE
	company.CreatedAt = time.Now().UTC()
	company.UpdatedAt = company.CreatedAt

	result, err := dbTransaction.ExecContext(
		ctx,
		`INSERT INTO companies(name,cnpj,status,created_at,updated_at) VALUES(?,?,?,?,?)`,
		company.Name,
		company.CNPJ,
		company.Status,
		company.CreatedAt,
		company.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to insert company")

		return model.Company{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Error to get last insert company id")

		return model.Company{}, err
	}

	company.ID = lastInsertID

	return company, nil
}
