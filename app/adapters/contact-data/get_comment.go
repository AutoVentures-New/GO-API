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

func GetComments(
	ctx context.Context,
	account string,
	ulids []string) ([]model.Comment, error) {

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListCommentsData, account) + " WHERE commented_at IN (" + whereIn + ") ORDER BY created_at ASC"

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list comments data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(
			&c.RecordNumber,
			&c.Ulid,
			&c.CreatedBy,
			&c.Text,
			&c.CommentedAt,
			&c.Files,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
