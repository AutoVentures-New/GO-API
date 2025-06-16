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

func GetActivityFiles(
	ctx context.Context,
	account string,
	ulids []string) ([]model.ActivityFile, error) {

	placeholders := make([]string, len(ulids))
	args := make([]interface{}, len(ulids))

	for i, ulid := range ulids {
		placeholders[i] = "?"
		args[i] = ulid
	}

	whereIn := strings.Join(placeholders, ", ")

	sqlQuery := fmt.Sprintf(query.ListActivityFileData, account) + " WHERE ulid IN (" + whereIn + ")"

	rows, err := database.Database.QueryContext(ctx, sqlQuery, args...)

	if err != nil {
		logrus.WithError(err).
			Error("Error to list activity file data")

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var results []model.ActivityFile
	for rows.Next() {
		var af model.ActivityFile
		err := rows.Scan(
			&af.RecordNumber,
			&af.Ulid,
			&af.CreatedBy,
			&af.To,
			&af.Subject,
			&af.Done,
			&af.Files,
			&af.UserCreateDate,
			&af.CreatedAt,
			&af.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, af)
	}

	comments, err := GetComments(ctx, account, pkg.ExtractIdentifiers(results))

	if err != nil {
		return nil, err
	}

	fileUlids := pkg.ExtractField(results, func(item model.ActivityFile) string {
		return item.CreatedBy
	})
	commentsUlid := pkg.ExtractField(comments, func(item model.Comment) string {
		return item.CreatedBy
	})

	users, err := GetUsers(ctx, append(fileUlids, commentsUlid...))

	if err != nil {
		return nil, err
	}

	usersMap := pkg.SliceToMap(users, func(c model.User) string { return c.Ulid })

	return GroupComments(results, comments, usersMap), nil
}

func GroupComments(files []model.ActivityFile, comments []model.Comment, users map[string]model.User) []model.ActivityFile {
	commentsMap := make(map[string][]model.Comment)

	setUserInfo := func(user model.User) (name, image string) {
		var parts []string
		if user.FirstName != nil && strings.TrimSpace(*user.FirstName) != "" {
			parts = append(parts, *user.FirstName)
		}
		if user.LastName != nil && strings.TrimSpace(*user.LastName) != "" {
			parts = append(parts, *user.LastName)
		}
		if user.Image != nil {
			image = *user.Image
		}
		return strings.Join(parts, " "), image
	}

	for i := range comments {
		if user, ok := users[comments[i].CreatedBy]; ok {
			name, image := setUserInfo(user)
			if name != "" {
				comments[i].CreatedByName = name
			}
			if image != "" {
				comments[i].CreatedByImage = image
			}
		}
		if comments[i].CommentedAt != nil {
			commentsMap[*comments[i].CommentedAt] = append(commentsMap[*comments[i].CommentedAt], comments[i])
		}
	}

	for i := range files {
		if user, ok := users[files[i].CreatedBy]; ok {
			name, image := setUserInfo(user)
			if name != "" {
				files[i].CreatedByName = name
			}
			if image != "" {
				files[i].CreatedByImage = image
			}
		}
		files[i].Comments = commentsMap[files[i].Ulid]
	}

	return files
}
