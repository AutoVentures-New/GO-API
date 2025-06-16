package contact_data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AutoVentures-New/GO-API/database"
	"github.com/AutoVentures-New/GO-API/internal/query"
	"github.com/AutoVentures-New/GO-API/model"
	"github.com/sirupsen/logrus"
	"strings"
)

func GetUsers(
	ctx context.Context,
	ulids []string,
) ([]model.User, error) {

	if len(ulids) == 0 {
		return []model.User{}, nil
	}

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListUsersData) + " AND ulid IN (" + whereIn + ")"

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list users data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	usersData := make([]model.User, 0)

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.Ulid,
			&user.FirstName,
			&user.LastName,
			&user.Image,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to list users data")
			return nil, err
		}
		usersData = append(usersData, user)
	}

	return usersData, nil
}
