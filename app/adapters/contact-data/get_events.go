package contact_data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AutoVentures-New/GO-API/database"
	"github.com/AutoVentures-New/GO-API/internal/query"
	"github.com/AutoVentures-New/GO-API/model"
	"github.com/AutoVentures-New/GO-API/pkg"
	"github.com/sirupsen/logrus"
	"strings"
)

func GetEvents(
	ctx context.Context,
	account string,
	ulids []string,
) ([]model.CalendarEvent, error) {

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListCalendarEventData, account, account) + " WHERE ce.ulid IN (" + whereIn + ")"

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
			&ev.CalendarUlid,
			&ev.Name,
			&ev.Description,
			&ev.Participants,
			&ev.When,
			&ev.Location,
			&ev.Recurrence,
			&ev.Notifications,
			&ev.Conferencing,
			&ev.Records,
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

	eventFiles, err := GetEventFiles(ctx, account, pkg.ExtractIdentifiers(events))

	if err != nil {
		return nil, err
	}

	calendars, err := GetCalendars(ctx, account)
	if err != nil {
		return nil, err
	}
	calendarMap := pkg.SliceToMap(calendars, func(c model.Calendar) string { return c.Ulid })

	return GroupFiles(events, eventFiles, calendarMap), nil
}

func GroupFiles(events []model.CalendarEvent, files []model.CalendarEventFile, calendarMap map[string]model.Calendar) []model.CalendarEvent {
	filesMap := make(map[string][]model.CalendarEventFile)

	for _, f := range files {
		if f.EventUlid != nil {
			filesMap[*f.EventUlid] = append(filesMap[*f.EventUlid], f)
		}
	}
	for i, e := range events {
		if calendar, ok := calendarMap[e.CalendarUlid]; ok {
			events[i].Calendar = calendar
		}

		if fileData, ok := filesMap[e.Ulid]; ok {
			events[i].ActivityFiles = fileData
			events[i].ConferenceRecords = []interface{}{}
		} else {
			events[i].ActivityFiles = []model.CalendarEventFile{}
			events[i].ConferenceRecords = []interface{}{}
		}
	}
	return events
}
