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

func GetCalls(
	ctx context.Context,
	account string,
	ulids []string,
	contact string,
) ([]model.Call, error) {

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListCallData, account) + " WHERE ulid IN (" + whereIn + ")"

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
			&call.Subject,
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

	callsUlid := pkg.ExtractField(callsData, func(item model.Call) string {
		return item.CreatedBy
	})

	users, err := GetUsers(ctx, callsUlid)

	if err != nil {
		return nil, err
	}

	usersMap := pkg.SliceToMap(users, func(c model.User) string { return c.Ulid })

	contacts, err := GetContacts(ctx, account, []string{contact})

	if err != nil {
		return nil, err
	}

	return GroupUserAndContacts(callsData, usersMap, contacts), nil

}

func GroupUserAndContacts(data []model.Call, users map[string]model.User, contacts []model.Contact) []model.Call {

	for i, call := range data {
		if user, ok := users[call.CreatedBy]; ok {
			var parts []string
			if user.FirstName != nil && strings.TrimSpace(*user.FirstName) != "" {
				parts = append(parts, *user.FirstName)
			}
			if user.LastName != nil && strings.TrimSpace(*user.LastName) != "" {
				parts = append(parts, *user.LastName)
			}
			if len(parts) > 0 {
				data[i].CreatedByName = strings.Join(parts, " ")
			}
			if user.Image != nil {
				data[i].CreatedByImage = *user.Image
			}
		}

		if len(contacts) >= 1 {
			var parts []string
			constact := contacts[0]
			if constact.FirstName != nil && strings.TrimSpace(*constact.FirstName) != "" {
				parts = append(parts, *constact.FirstName)
			}
			if constact.LastName != nil && strings.TrimSpace(*constact.LastName) != "" {
				parts = append(parts, *constact.LastName)
			}
			if constact.CompanyName != nil && strings.TrimSpace(*constact.CompanyName) != "" {
				parts = append(parts, *constact.CompanyName)
			}
			if len(parts) > 0 {
				data[i].ContactName = strings.Join(parts, " ")
			}
			if constact.Image != nil {
				data[i].ContactImage = *constact.Image
			}

		}
	}
	return data
}
