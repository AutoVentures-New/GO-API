package profile

import (
	"context"
	"encoding/json"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

func GetCurriculum(
	ctx context.Context,
	candidateID int64,
) (model.Curriculum, error) {
	curriculum := model.Curriculum{}

	var languagesJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,candidate_id,gender,gender_identifier,color,is_special_needs,languages,created_at,updated_at 
				FROM candidate_curriculum WHERE candidate_id = ?`,
		candidateID,
	).Scan(
		&curriculum.ID,
		&curriculum.CandidateID,
		&curriculum.Gender,
		&curriculum.GenderIdentifier,
		&curriculum.Color,
		&curriculum.IsSpecialNeeds,
		&languagesJSON,
		&curriculum.CreatedAt,
		&curriculum.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get candidate curriculum")

		return curriculum, err
	}

	err = json.Unmarshal(languagesJSON, &curriculum.Languages)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal candidate curriculum languages")

		return curriculum, err
	}

	curriculum.ProfessionalExperiences, err = getCurriculumProfessionalExperience(ctx, candidateID)
	if err != nil {
		return curriculum, err
	}

	curriculum.AcademicExperiences, err = getCurriculumAcademicExperience(ctx, candidateID)
	if err != nil {
		return curriculum, err
	}

	return curriculum, nil
}

func getCurriculumAcademicExperience(
	ctx context.Context,
	candidateID int64,
) ([]model.AcademicExperience, error) {
	academicExperiences := make([]model.AcademicExperience, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id,candidate_id,title,company,area_id,status,level,start_date,end_date,created_at,updated_at 
				FROM candidate_curriculum_academic_experience WHERE candidate_id = ?`,
		candidateID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get candidate curriculum academic experience")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var academicExperience model.AcademicExperience

		err = rows.Scan(
			&academicExperience.ID,
			&academicExperience.CandidateID,
			&academicExperience.Title,
			&academicExperience.Company,
			&academicExperience.AreaID,
			&academicExperience.Status,
			&academicExperience.Level,
			&academicExperience.StartDate,
			&academicExperience.EndDate,
			&academicExperience.CreatedAt,
			&academicExperience.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to get candidate curriculum academic experience")

			return nil, err
		}

		academicExperiences = append(academicExperiences, academicExperience)
	}

	return academicExperiences, nil
}

func getCurriculumProfessionalExperience(
	ctx context.Context,
	candidateID int64,
) ([]model.ProfessionalExperience, error) {
	professionalExperiences := make([]model.ProfessionalExperience, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT id,candidate_id,title,company,area_id,city,state,job_mode,current_job,start_date,end_date,created_at,updated_at 
				FROM candidate_curriculum_professional_experience WHERE candidate_id = ?`,
		candidateID,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get candidate curriculum professional experience")

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var professionalExperience model.ProfessionalExperience

		err = rows.Scan(
			&professionalExperience.ID,
			&professionalExperience.CandidateID,
			&professionalExperience.Title,
			&professionalExperience.Company,
			&professionalExperience.AreaID,
			&professionalExperience.City,
			&professionalExperience.State,
			&professionalExperience.JobMode,
			&professionalExperience.CurrentJob,
			&professionalExperience.StartDate,
			&professionalExperience.EndDate,
			&professionalExperience.CreatedAt,
			&professionalExperience.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to get candidate curriculum professional experience")

			return nil, err
		}

		professionalExperiences = append(professionalExperiences, professionalExperience)
	}

	return professionalExperiences, nil
}
