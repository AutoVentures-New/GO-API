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

func GetEmails(
	ctx context.Context,
	user model.User,
	ulids []string,
) ([]model.EmailBucket, error) {

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListEmailData, user.Account) + " WHERE ulid IN (" + whereIn + ")"

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list email data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	emailsData := make([]model.EmailBucket, 0)

	for rows.Next() {
		var e model.EmailBucket
		err := rows.Scan(
			&e.RecordNumber,
			&e.Ulid,
			&e.AccountID,
			&e.MessageID,
			&e.ThreadID,
			&e.Subject,
			&e.From,
			&e.To,
			&e.CC,
			&e.BCC,
			&e.ReplyTo,
			&e.Headers,
			&e.Starred,
			&e.Unread,
			&e.ReplyToMessageID,
			&e.Body,
			&e.Files,
			&e.Folder,
			&e.Links,
			&e.Opens,
			&e.LinkClicks,
			&e.IsTracked,
			&e.Date,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to list email data")
			return nil, err
		}
		emailsData = append(emailsData, e)
	}

	return emailsData, nil
}
