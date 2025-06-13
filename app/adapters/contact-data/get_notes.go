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

func GetNotes(
	ctx context.Context,
	account string,
	ulids []string) ([]model.Note, error) {

	notes, err := GetNotesSelect(ctx, account, ulids, false)

	if err != nil {
		return nil, err
	}

	notesComments, err := GetNotesSelect(ctx, account, pkg.ExtractIdentifiers(notes), true)

	return GroupNotes(notes, notesComments), nil
}

func GetNotesSelect(
	ctx context.Context,
	account string,
	ulids []string,
	isComments bool,
) ([]model.Note, error) {

	if len(ulids) == 0 {
		return []model.Note{}, nil
	}

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	condition := "ulid"
	if isComments {
		condition = "commented_at"
	}

	sqlQuery := fmt.Sprintf(
		`%s WHERE %s IN (%s) ORDER BY n.created_at ASC`,
		fmt.Sprintf(query.ListNoteData, account),
		condition,
		whereIn,
	)
	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list note data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var notes []model.Note
	for rows.Next() {
		var note model.Note
		err := rows.Scan(
			&note.RecordNumber,
			&note.Ulid,
			&note.CreatedBy,
			&note.To,
			&note.Subject,
			&note.Done,
			&note.Text,
			&note.CommentedAt,
			&note.Files,
			&note.UserCreateDate,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func GroupNotes(notes []model.Note, comments []model.Note) []model.Note {
	commentsMap := make(map[string][]model.Note)

	for _, c := range comments {
		if c.CommentedAt != nil {
			commentsMap[*c.CommentedAt] = append(commentsMap[*c.CommentedAt], c)
		}
	}
	for i, n := range notes {
		if replies, ok := commentsMap[n.Ulid]; ok {
			notes[i].Comments = replies
		} else {
			notes[i].Comments = []model.Note{}
		}
	}
	return notes
}
