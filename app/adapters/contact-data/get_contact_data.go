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
) ([]model.ContactData, int, error) {

	sqlQuery := fmt.Sprintf(query.ListContactData, account, account)
	countQuery := fmt.Sprintf(query.CountContactData, account, account)

	var args []interface{}

	var total = 0

	whereClause, args := getWhereClause(filter, args)

	sqlQuery += whereClause

	countQuery += whereClause

	sqlQuery = getOrderByClause(filter, sqlQuery)

	limit := filter.Limit
	offset := filter.Page * limit

	sqlQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	err := database.Database.QueryRowContext(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {
		return nil, total, err

	}

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list providers")

		return nil, total, err
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

			return nil, total, err
		}

		contactsData = append(contactsData, contactData)
	}

	return contactsData, total, nil
}

func getWhereClause(filter request.ContactDataQuery, args []interface{}) (string, []interface{}) {
	sqlQuery := ""
	if filter.ContactULID != nil {
		sqlQuery += "AND cdc.contact_ulid = ?"
		args = append(args, filter.ContactULID)
	}

	if filter.Type != "" && filter.Type != model.TypeAll && filter.OrderBy != nil && *filter.OrderBy != "new_emails" && filter.Type != model.TypeCall && filter.Type != model.TypePhoneCall {
		sqlQuery += " AND cd.type = ?"
		args = append(args, filter.Type)
	}

	if filter.Type == model.TypeCall || filter.Type == model.TypePhoneCall {
		sqlQuery += " AND (cd.type = ? OR cd.type = ?)"
		args = append(args, "CALL", "PHONE_CALL")
	}

	if filter.OrderBy != nil && *filter.OrderBy == "new_emails" {
		sqlQuery += " AND cd.type = ? "
		args = append(args, "EMAIL")

		sqlQuery += " AND cdc.is_new = ?"
		args = append(args, true)
	}

	if filter.Unread != nil && filter.OrderBy != nil && *filter.OrderBy == "new_emails" || filter.Type == model.TypeEmail {
		sqlQuery += " AND cd.unread = ?"
		value := 0
		if *filter.Unread {
			value = 1
		}

		args = append(args, value)
	}

	if filter.DateFrom != nil && filter.DateTo != nil {
		sqlQuery += " AND cd.date >= ? AND cd.date <= ?"
		args = append(args, filter.DateFrom, filter.DateTo)
	}

	if filter.From != nil {
		sqlQuery += " AND cd.from LIKE ?"
		args = append(args, "%"+*filter.From+"%")
	}

	if filter.To != nil {
		sqlQuery += " AND cd.to LIKE ?"
		args = append(args, "%"+*filter.To+"%")
	}

	if filter.Subject != nil {
		sqlQuery += " AND cd.search_field LIKE ?"
		args = append(args, "%"+*filter.Subject+"%")
	}

	if filter.Folder != nil {
		sqlQuery += " AND cd.folder = ?"
		args = append(args, filter.Folder)
	}

	return sqlQuery, args
}

func getOrderByClause(filter request.ContactDataQuery, sqlQuery string) string {
	var column string

	if filter.OrderBy == nil || *filter.OrderBy == "created_at" {
		column = "cd.created_at"
	} else {
		switch *filter.OrderBy {
		case "date", "updated_at":
			column = "cd." + *filter.OrderBy
		case "subject", "title":
			column = "cd.search_field"
		default:
			column = "cd.created_at"
		}
	}

	sort := "DESC"
	if filter.Sort != nil && (*filter.Sort == "ASC" || *filter.Sort == "DESC") {
		sort = *filter.Sort
	}

	sqlQuery += fmt.Sprintf(" ORDER BY %s %s ", column, sort)
	return sqlQuery
}
