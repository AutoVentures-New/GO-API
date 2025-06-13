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

func GetActivityFiles(
	ctx context.Context,
	user model.User,
	ulids []string) ([]model.ActivityFile, error) {

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListActivityFileData, user.Account) + " WHERE ulid IN (" + whereIn + ")"

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list activity file data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var results []model.ActivityFile
	for rows.Next() {
		var af model.ActivityFile
		err := rows.Scan(
			&af.RecordNumber,
			&af.Ulid,
			&af.CreatedBy,
			&af.To,
			&af.Subject,
			&af.Done,
			&af.Files,
			&af.UserCreateDate,
			&af.CreatedAt,
			&af.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, af)
	}

	return results, nil
}
