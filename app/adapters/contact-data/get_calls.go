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

func GetCalls(
	ctx context.Context,
	user model.User,
	ulids []string,
) ([]model.Call, error) {

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListCallData, user.Account) + " WHERE ulid IN (" + whereIn + ")"

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list call data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	callsData := make([]model.Call, 0)

	for rows.Next() {
		var call model.Call
		err := rows.Scan(
			&call.RecordNumber,
			&call.Ulid,
			&call.CreatedBy,
			&call.UserPhoneNumber,
			&call.ContactPhoneNumber,
			&call.Done,
			&call.To,
			&call.CallType,
			&call.Direction,
			&call.Outcome,
			&call.Notes,
			&call.UserCreateDate,
			&call.CreatedAt,
			&call.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to list call data")
			return nil, err
		}
		callsData = append(callsData, call)
	}

	return callsData, nil
}
