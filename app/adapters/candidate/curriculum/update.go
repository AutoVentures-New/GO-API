package profile

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func UpdateCurriculum(
	ctx context.Context,
	curriculum model.Curriculum,
) (model.Curriculum, error) {
	dbTransaction, err := database.Database.Begin()
	if err != nil {
		logrus.WithError(err).Error("Error to open db transaction")

		return model.Curriculum{}, err
	}

	languagesString, err := json.Marshal(curriculum.Languages)
	if err != nil {
		logrus.WithError(err).Error("Error to marshal candidate languages")

		_ = dbTransaction.Rollback()

		return model.Curriculum{}, err
	}

	curriculum.UpdatedAt = time.Now().UTC()

	_, err = dbTransaction.ExecContext(
		ctx,
		`UPDATE candidate_curriculum set
    				gender = ?, 
    				is_special_needs = ?,  
    				languages = ?,  
    				updated_at = ?
				WHERE candidate_id = ?`,
		curriculum.Gender,
		curriculum.IsSpecialNeeds,
		languagesString,
		curriculum.UpdatedAt,
		curriculum.CandidateID,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to update candidate curriculum")

		return model.Curriculum{}, err
	}

	curriculum.ProfessionalExperiences, err = updateCurriculumProfessionalExperience(
		ctx,
		curriculum,
		dbTransaction,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		return model.Curriculum{}, err
	}

	curriculum.AcademicExperiences, err = updateCurriculumAcademicExperience(
		ctx,
		curriculum,
		dbTransaction,
	)
	if err != nil {
		_ = dbTransaction.Rollback()

		return model.Curriculum{}, err
	}

	err = dbTransaction.Commit()
	if err != nil {
		_ = dbTransaction.Rollback()

		logrus.WithError(err).Error("Error to commit db transaction")

		return model.Curriculum{}, err
	}

	return curriculum, nil
}

func updateCurriculumProfessionalExperience(
	ctx context.Context,
	curriculum model.Curriculum,
	dbTransaction *sql.Tx,
) ([]model.ProfessionalExperience, error) {
	newExperience := make([]model.ProfessionalExperience, 0)

	ignoreIds := make([]string, 0)

	for _, item := range curriculum.ProfessionalExperiences {
		if item.ID > 0 {
			item.UpdatedAt = time.Now().UTC()
			_, err := dbTransaction.ExecContext(
				ctx,
				`UPDATE candidate_curriculum_professional_experience set
    				title = ?,
					company = ?,
					area_id = ?,
					city = ?,
					state = ?,
					job_mode = ?,
					current_job = ?,
					start_date = ?,
					end_date = ?,
    				updated_at = ?
				WHERE candidate_id = ? AND id = ?`,
				item.Title,
				item.Company,
				item.AreaID,
				item.City,
				item.State,
				item.JobMode,
				item.CurrentJob,
				item.StartDate,
				item.EndDate,
				item.UpdatedAt,
				curriculum.CandidateID,
				item.ID,
			)
			if err != nil {
				logrus.WithError(err).Error("Error to update candidate curriculum_professional_experience")

				return newExperience, err
			}

			newExperience = append(newExperience, item)
			ignoreIds = append(ignoreIds, fmt.Sprintf("%d", item.ID))

			continue
		}

		item.CandidateID = curriculum.CandidateID
		item.CreatedAt = time.Now().UTC()
		item.UpdatedAt = item.CreatedAt

		resultInsert, err := dbTransaction.ExecContext(
			ctx,
			`INSERT INTO candidate_curriculum_professional_experience(candidate_id,title,company,area_id,city,state,job_mode,current_job,start_date,end_date,created_at,updated_at)
					VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`,
			item.CandidateID,
			item.Title,
			item.Company,
			item.AreaID,
			item.City,
			item.State,
			item.JobMode,
			item.CurrentJob,
			item.StartDate,
			item.EndDate,
			item.CreatedAt,
			item.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to insert candidate curriculum_professional_experience")

			return newExperience, err
		}

		lastId, err := resultInsert.LastInsertId()
		if err != nil {
			logrus.WithError(err).Error("Error to get insert id candidate curriculum_professional_experience")

			return newExperience, err
		}

		item.ID = lastId
		newExperience = append(newExperience, item)
		ignoreIds = append(ignoreIds, fmt.Sprintf("%d", item.ID))
	}

	query := fmt.Sprintf(
		`DELETE FROM candidate_curriculum_professional_experience WHERE candidate_id = ? AND id NOT IN(%s)`,
		strings.Join(ignoreIds, ","),
	)

	if len(ignoreIds) == 0 {
		query = `DELETE FROM candidate_curriculum_professional_experience WHERE candidate_id = ?`
	}

	_, err := dbTransaction.ExecContext(
		ctx,
		query,
		curriculum.CandidateID,
	)
	if err != nil {
		logrus.WithError(err).
			WithField("query", query).
			Error("Error to delete candidate curriculum_professional_experience")

		return nil, err
	}

	return newExperience, nil
}

func updateCurriculumAcademicExperience(
	ctx context.Context,
	curriculum model.Curriculum,
	dbTransaction *sql.Tx,
) ([]model.AcademicExperience, error) {
	newExperience := make([]model.AcademicExperience, 0)

	ignoreIds := make([]string, 0)

	for _, item := range curriculum.AcademicExperiences {
		if item.ID > 0 {
			item.UpdatedAt = time.Now().UTC()
			_, err := dbTransaction.ExecContext(
				ctx,
				`UPDATE candidate_curriculum_academic_experience set
    				title = ?,
					company = ?,
					area_id = ?,
					status = ?,
					level = ?,
					start_date = ?,
					end_date = ?,
    				updated_at = ?
				WHERE candidate_id = ? AND id = ?`,
				item.Title,
				item.Company,
				item.AreaID,
				item.Status,
				item.Level,
				item.StartDate,
				item.EndDate,
				item.UpdatedAt,
				curriculum.CandidateID,
				item.ID,
			)
			if err != nil {
				logrus.WithError(err).Error("Error to update candidate curriculum_academic_experience")

				return newExperience, err
			}

			newExperience = append(newExperience, item)
			ignoreIds = append(ignoreIds, fmt.Sprintf("%d", item.ID))

			continue
		}

		item.CandidateID = curriculum.CandidateID
		item.CreatedAt = time.Now().UTC()
		item.UpdatedAt = item.CreatedAt

		resultInsert, err := dbTransaction.ExecContext(
			ctx,
			`INSERT INTO candidate_curriculum_academic_experience(candidate_id,title,company,area_id,status,level,start_date,end_date,created_at,updated_at)
					VALUES(?,?,?,?,?,?,?,?,?,?)`,
			item.CandidateID,
			item.Title,
			item.Company,
			item.AreaID,
			item.Status,
			item.Level,
			item.StartDate,
			item.EndDate,
			item.CreatedAt,
			item.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to insert candidate curriculum_academic_experience")

			return newExperience, err
		}

		lastId, err := resultInsert.LastInsertId()
		if err != nil {
			logrus.WithError(err).Error("Error to get insert id candidate curriculum_academic_experience")

			return newExperience, err
		}

		item.ID = lastId
		newExperience = append(newExperience, item)
		ignoreIds = append(ignoreIds, fmt.Sprintf("%d", item.ID))
	}

	query := fmt.Sprintf(
		`DELETE FROM candidate_curriculum_academic_experience WHERE candidate_id = ? AND id NOT IN(%s)`,
		strings.Join(ignoreIds, ","),
	)

	if len(ignoreIds) == 0 {
		query = `DELETE FROM candidate_curriculum_academic_experience WHERE candidate_id = ?`
	}

	_, err := dbTransaction.ExecContext(
		ctx,
		query,
		curriculum.CandidateID,
	)
	if err != nil {
		logrus.WithError(err).
			WithField("query", query).
			Error("Error to delete candidate curriculum_academic_experience")

		return nil, err
	}

	return newExperience, nil
}
