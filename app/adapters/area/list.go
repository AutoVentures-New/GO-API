package area

import (
	"context"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListAreas(
	ctx context.Context,
) ([]model.Area, error) {
	areas := make([]model.Area, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id,title,created_at,updated_at FROM areas`,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list areas")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		area := model.Area{}
		err := rows.Scan(
			&area.ID,
			&area.Title,
			&area.CreatedAt,
			&area.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan area")

			return nil, err
		}

		areas = append(areas, area)
	}

	return areas, nil
}
