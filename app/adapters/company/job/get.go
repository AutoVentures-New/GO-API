package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrJobNotFound = errors.New("Job not found")

func GetJob(
	ctx context.Context,
	id int64,
	companyID int64,
) (model.Job, error) {
	job := model.Job{}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,title,company_id,area_id,is_talent_bank,is_special_needs,description,job_mode,contracting_modality,state,city,responsibilities,questionnaire,video_link,status,publish_at,finish_at,created_at,updated_at 
				FROM jobs WHERE company_id = ? AND id = ?`,
		companyID,
		id,
	).Scan(
		&job.ID,
		&job.Title,
		&job.CompanyID,
		&job.AreaID,
		&job.IsTalentBank,
		&job.IsSpecialNeeds,
		&job.Description,
		&job.JobMode,
		&job.ContractingModality,
		&job.State,
		&job.City,
		&job.Responsibilities,
		&job.Questionnaire,
		&job.VideoLink,
		&job.Status,
		&job.PublishAt,
		&job.FinishAt,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return job, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job")

		return job, err
	}

	job.JobCulturalFit, err = getJobCulturalFit(ctx, id, companyID)
	if err != nil {
		return job, err
	}

	job.JobRequirement, err = getJobRequirements(ctx, id, companyID)
	if err != nil {
		return job, err
	}

	job.Benefits, err = getJobBenefits(ctx, id, companyID)
	if err != nil {
		return job, err
	}

	job.VideoQuestions, err = getJobVideoQuestions(ctx, id, companyID)
	if err != nil {
		return job, err
	}

	job.Questions, err = getJobQuestions(ctx, id, companyID)
	if err != nil {
		return job, err
	}

	job.Area, err = getArea(ctx, job.AreaID)
	if err != nil {
		return job, err
	}

	return job, nil
}

func getJobCulturalFit(
	ctx context.Context,
	id int64,
	companyID int64,
) (*model.JobCulturalFit, error) {
	jobCulturalFit := model.JobCulturalFit{}

	var answersJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,company_id,job_id,answers,created_at,updated_at 
				FROM job_cultural_fit WHERE company_id = ? AND job_id = ? LIMIT 1`,
		companyID,
		id,
	).Scan(
		&jobCulturalFit.ID,
		&jobCulturalFit.CompanyID,
		&jobCulturalFit.JobID,
		&answersJSON,
		&jobCulturalFit.CreatedAt,
		&jobCulturalFit.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job cultural fit")

		return nil, err
	}

	err = json.Unmarshal(answersJSON, &jobCulturalFit.Answers)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal job cultural fit")

		return nil, err
	}

	return &jobCulturalFit, nil
}

func getJobRequirements(
	ctx context.Context,
	id int64,
	companyID int64,
) (*model.JobRequirement, error) {
	jobRequirement := model.JobRequirement{}

	var itemsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,company_id,job_id,items,min_match,created_at,updated_at 
				FROM job_requirements WHERE company_id = ? AND job_id = ?`,
		companyID,
		id,
	).Scan(
		&jobRequirement.ID,
		&jobRequirement.CompanyID,
		&jobRequirement.JobID,
		&itemsJSON,
		&jobRequirement.MinMatch,
		&jobRequirement.CreatedAt,
		&jobRequirement.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job requirement")

		return nil, err
	}

	if err = json.Unmarshal(itemsJSON, &jobRequirement.Items); err != nil {
		logrus.WithError(err).Error("Error to unmarshal job requirement")

		return nil, err
	}

	return &jobRequirement, nil
}

func getJobBenefits(
	ctx context.Context,
	id int64,
	companyID int64,
) ([]model.Benefit, error) {
	benefits := make([]model.Benefit, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT benefits.id, benefits.name, benefits.company_id, benefits.created_at, benefits.updated_at FROM job_benefits
				JOIN benefits ON job_benefits.benefit_id = benefits.id 
				WHERE job_benefits.company_id = ? AND job_benefits.job_id = ?`,
		companyID,
		id,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get job benefits")

		return benefits, err
	}

	defer rows.Close()

	for rows.Next() {
		benefit := model.Benefit{}

		err := rows.Scan(
			&benefit.ID,
			&benefit.Name,
			&benefit.CompanyID,
			&benefit.CreatedAt,
			&benefit.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan job benefit")

			return benefits, err
		}

		benefits = append(benefits, benefit)
	}

	return benefits, nil
}

func getJobVideoQuestions(
	ctx context.Context,
	id int64,
	companyID int64,
) (*model.JobVideoQuestions, error) {
	jobVideoQuestions := model.JobVideoQuestions{}

	var questionsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,company_id,job_id,questions,created_at,updated_at 
				FROM job_video_questions WHERE company_id = ? AND job_id = ?`,
		companyID,
		id,
	).Scan(
		&jobVideoQuestions.ID,
		&jobVideoQuestions.CompanyID,
		&jobVideoQuestions.JobID,
		&questionsJSON,
		&jobVideoQuestions.CreatedAt,
		&jobVideoQuestions.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job video questions")

		return nil, err
	}

	if err = json.Unmarshal(questionsJSON, &jobVideoQuestions.Questions); err != nil {
		logrus.WithError(err).Error("Error to unmarshal job video questions")

		return nil, err
	}

	return &jobVideoQuestions, nil
}

func getJobQuestions(
	ctx context.Context,
	id int64,
	companyID int64,
) ([]model.Question, error) {
	questions := make([]model.Question, 0)

	rows, err := database.Database.QueryContext(
		ctx,
		`SELECT 
    				questionnaire_questions.id, 
    				questionnaire_questions.title, 
    				questionnaire_questions.type, 
    				questionnaire_questions.questionnaire_id, 
    				questionnaire_questions.created_at, 
    				questionnaire_questions.updated_at 
				FROM job_questions
				JOIN questionnaire_questions ON job_questions.question_id = questionnaire_questions.id 
				WHERE job_questions.company_id = ? AND job_questions.job_id = ?`,
		companyID,
		id,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get job questions")

		return questions, err
	}

	defer rows.Close()

	for rows.Next() {
		question := model.Question{}

		err := rows.Scan(
			&question.ID,
			&question.Title,
			&question.Type,
			&question.QuestionnaireID,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Error to scan job questions")

			return questions, err
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func getArea(
	ctx context.Context,
	id int64,
) (*model.Area, error) {
	area := model.Area{}

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,title,created_at,updated_at 
				FROM areas WHERE id = ?`,
		id,
	).Scan(
		&area.ID,
		&area.Title,
		&area.CreatedAt,
		&area.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrJobNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get area")

		return nil, err
	}

	return &area, nil
}
