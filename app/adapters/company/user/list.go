package user

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListUsers(ctx context.Context, companyID int64, userLogged int64) ([]model.User, error) {
	users := make([]model.User, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id,name,cpf,phone,role,email,status,company_id,created_at,updated_at 
				FROM users 
				WHERE company_id = ? and id != ?`,
		companyID,
		userLogged,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list users")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.CPF,
			&user.Phone,
			&user.Role,
			&user.Email,
			&user.Status,
			&user.CompanyID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan users")

			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
