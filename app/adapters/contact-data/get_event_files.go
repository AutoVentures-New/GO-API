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

func GetEventFiles(ctx context.Context, account string, ulids []string,
) ([]model.CalendarEventFile, error) {
	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListEventFileData, account) + " WHERE event_ulid IN (" + whereIn + ") ORDER BY created_at DESC"

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list event file data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var files []model.CalendarEventFile

	for rows.Next() {
		var f model.CalendarEventFile
		err := rows.Scan(
			&f.RecordNumber,
			&f.Ulid,
			&f.EventUlid,
			&f.IsExternal,
			&f.Name,
			&f.Extension,
			&f.Link,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}
