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

func GetContacts(
	ctx context.Context,
	account string,
	ulids []string,
) ([]model.Contact, error) {

	if len(ulids) == 0 {
		return []model.Contact{}, nil
	}

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListContactsData, account, account) + " AND ulid IN (" + whereIn + ")"

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list contacts data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	contactsData := make([]model.Contact, 0)

	for rows.Next() {
		var contact model.Contact
		err := rows.Scan(
			&contact.Ulid,
			&contact.FirstName,
			&contact.LastName,
			&contact.CompanyName,
			&contact.Image,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to list contacts data")
			return nil, err
		}
		contactsData = append(contactsData, contact)
	}

	return contactsData, nil
}
