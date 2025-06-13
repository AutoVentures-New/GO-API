package contact_data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AutoVentures-New/GO-API/database"
	"github.com/AutoVentures-New/GO-API/model"
	"github.com/AutoVentures-New/GO-API/model/request"
	"github.com/sirupsen/logrus"
)

func GetContactData(
	ctx context.Context,
	user model.User,
	filter request.ContactDataQuery,
) ([]model.ContactData, error) {

	listContactDataQuery := fmt.Sprintf("SELECT cd.ulid, cd.date FROM tenant_%s.contact_data_contact cdc "+
		"INNER JOIN tenant_%s.contact_data cd "+
		"ON cdc.contact_data_ulid COLLATE utf8mb4_unicode_ci = cd.ulid COLLATE utf8mb4_unicode_ci "+
		"WHERE cdc.contact_ulid = ?", user.Account, user.Account)

	rows, err := database.Database.QueryContext(ctx, listContactDataQuery, filter.ContactULID)

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
