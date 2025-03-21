package user

import (
	"context"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateUser(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	user.UpdatedAt = time.Now().UTC()

	_, err := database.Database.ExecContext(
		ctx,
		`UPDATE users set name=?,email=?,cpf=?,status=?,updated_at = ? WHERE id = ? AND company_id = ?`,
		user.Name,
		user.Email,
		user.CPF,
		user.Status,
		user.UpdatedAt,
		user.ID,
		user.CompanyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to update user")

		return user, err
	}

	return user, nil
}
