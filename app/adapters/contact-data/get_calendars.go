package contact_data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AutoVentures-New/GO-API/database"
	"github.com/AutoVentures-New/GO-API/internal/query"
	"github.com/AutoVentures-New/GO-API/model"
	"github.com/sirupsen/logrus"
)

func GetCalendars(
	ctx context.Context,
	account string,
) ([]model.Calendar, error) {

	sqlQuery := fmt.Sprintf(query.ListCalendarData, account)

	rows, err := database.Database.QueryContext(ctx, sqlQuery)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list calendar data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	calendarData := make([]model.Calendar, 0)

	for rows.Next() {
		var c model.Calendar
		err := rows.Scan(
			&c.Ulid,
			&c.UserUlid,
			&c.ProviderUlid,
			&c.ExternalID,
			&c.Name,
			&c.Color,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to list calendar data")
			return nil, err
		}
		calendarData = append(calendarData, c)
	}

	return calendarData, nil
}
