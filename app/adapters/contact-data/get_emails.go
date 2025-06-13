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
	account string,
	ulids []string,
) ([]model.EmailBucket, error) {

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListEmailData, account) + " WHERE thread_id IN (" + whereIn + ") ORDER  BY eb.date DESC"

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

	return GroupEmailsByThreadID(emailsData), nil
}

func GroupEmailsByThreadID(emails []model.EmailBucket) []model.EmailBucket {
	grouped := make([]model.EmailBucket, 0)
	processed := make(map[string]bool)
	threads := make(map[string][]model.EmailBucket)

	for _, e := range emails {
		threads[e.ThreadID] = append(threads[e.ThreadID], e)
	}

	for _, group := range threads {
		if len(group) == 0 {
			continue
		}
		parent := group[0]
		parent.EmailsReply = group[1:]
		grouped = append(grouped, parent)

		for _, e := range group {
			processed[e.Ulid] = true
		}
	}

	for _, e := range emails {
		if !processed[e.Ulid] {
			grouped = append(grouped, e)
		}
	}

	return grouped
}
