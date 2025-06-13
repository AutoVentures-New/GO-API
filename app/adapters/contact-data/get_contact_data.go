package contact_data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AutoVentures-New/GO-API/database"
	"github.com/AutoVentures-New/GO-API/internal/query"
	"github.com/AutoVentures-New/GO-API/model"
	"github.com/AutoVentures-New/GO-API/model/request"
	"github.com/sirupsen/logrus"
)

func GetContactData(
	ctx context.Context,
	account string,
	filter request.ContactDataQuery,
) ([]model.ContactData, error) {

	sqlQuery := fmt.Sprintf(query.ListContactData, account, account)

	rows, err := database.Database.QueryContext(ctx, sqlQuery, filter.ContactULID)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list providers")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	contactsData := make([]model.ContactData, 0)

	for rows.Next() {
		var contactData model.ContactData

		err := rows.Scan(
			&contactData.Ulid,
			&contactData.Type,
			&contactData.Identifier,
			&contactData.From,
			&contactData.To,
			&contactData.CC,
			&contactData.Date,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to list contact data")

			return nil, err
		}

		contactsData = append(contactsData, contactData)
	}

	return contactsData, nil
}
