package job

import (
	"context"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func ListStateCities(
	ctx context.Context,
	companyID int64,
) (map[string]model.State, error) {
	states := make(map[string]model.State)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT state,city FROM jobs WHERE company_id = ? and status = 'ACTIVE'`,
		companyID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to list states and cities")

		return nil, err
	}

	defer rows.Close()

	cache := make(map[string]bool)

	for rows.Next() {
		var state, city string

		err := rows.Scan(
			&state,
			&city,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan state city")

			return nil, err
		}

		if _, ok := cache[state+city]; ok {
			continue
		}

		cache[state+city] = true

		stateModel := model.State{
			Name:   state,
			Cities: make([]model.City, 0),
		}

		if _, ok := states[state]; ok {
			stateModel = states[state]
		}

		stateModel.Cities = append(stateModel.Cities, model.City{
			Name: city,
		})

		states[state] = stateModel
	}

	return states, nil
}
