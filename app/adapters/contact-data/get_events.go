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

func GetEvents(
	ctx context.Context,
	user model.User,
	ulids []string,
) ([]model.CalendarEvent, error) {

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListCalendarEventData, user.Account) + " WHERE ulid IN (" + whereIn + ")"

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list event data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var events []model.CalendarEvent
	for rows.Next() {
		var ev model.CalendarEvent
		err := rows.Scan(
			&ev.RecordNumber,
			&ev.Ulid,
			&ev.Name,
			&ev.Description,
			&ev.Participants,
			&ev.When,
			&ev.Location,
			&ev.Recurrence,
			&ev.Notifications,
			&ev.Conferencing,
			&ev.ConferenceRecords,
			&ev.OrganizerName,
			&ev.OrganizerEmail,
			&ev.Owner,
			&ev.Done,
			&ev.StartDate,
			&ev.EndDate,
			&ev.AllDay,
			&ev.Type,
			&ev.Sequence,
			&ev.Files,
			&ev.CreatedAt,
			&ev.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to list event data")
			return nil, err
		}
		events = append(events, ev)
	}

	return events, nil
}
