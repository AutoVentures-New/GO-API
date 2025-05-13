package job

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	profile "github.com/hubjob/api/app/adapters/candidate/curriculum"
	questions_adp "github.com/hubjob/api/app/adapters/questions"
	"github.com/hubjob/api/database"
	"github.com/hubjob/api/model"
	"github.com/sirupsen/logrus"
)

var ErrApplicationNotFound = errors.New("Application not found")

func GetApplication(
	ctx context.Context,
	companyID int64,
	jobId int64,
	id int64,
) (model.Application, error) {
	application := model.Application{}

	var stepsJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,company_id,job_id,candidate_id,steps,current_step,status,created_at,updated_at 
				FROM job_applications WHERE company_id = ? AND job_id = ? AND id = ?`,
		companyID,
		jobId,
		id,
	).Scan(
		&application.ID,
		&application.CompanyID,
		&application.JobID,
		&application.CandidateID,
		&stepsJSON,
		&application.CurrentStep,
		&application.Status,
		&application.CreatedAt,
		&application.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return application, ErrApplicationNotFound
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job application")

		return application, err
	}

	err = json.Unmarshal(stepsJSON, &application.Steps)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal application steps")

		return application, err
	}

	application.Candidate, err = getCandidate(ctx, application.CandidateID)
	if err != nil {
		return application, err
	}

	application.JobApplicationRequirement, err = getJobApplicationRequirement(ctx, application.ID)
	if err != nil {
		return application, err
	}

	application.CulturalFit, err = getJobApplicationCulturalFit(ctx, application.ID, application.JobID)
	if err != nil {
		return application, err
	}

	application.JobApplicationQuestion, err = getJobApplicationQuestion(ctx, application.ID, application.JobID)
	if err != nil {
		return application, err
	}

	application.JobApplicationCandidateVideo, err = getJobApplicationCandidateVideo(ctx, application.ID)
	if err != nil {
		return application, err
	}

	return application, nil
}

func getCandidate(
	ctx context.Context,
	candidateID int64,
) (*model.Candidate, error) {
	var candidate model.Candidate

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT id,name,cpf,email,password,status,phone,birth_date,created_at,updated_at FROM candidates WHERE id = ?`,
		candidateID,
	).Scan(
		&candidate.ID,
		&candidate.Name,
		&candidate.CPF,
		&candidate.Email,
		&candidate.Password,
		&candidate.Status,
		&candidate.Phone,
		&candidate.BirthDate,
		&candidate.CreatedAt,
		&candidate.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Error to get candidate")

		return nil, err
	}

	_ = database.Database.QueryRowContext(
		ctx,
		`SELECT state, city FROM candidate_addresses WHERE candidate_id = ?`,
		candidateID,
	).Scan(
		&candidate.Address.State,
		&candidate.Address.City,
	)

	candidate.Curriculum, err = profile.GetCurriculum(ctx, candidate.ID)
	if err != nil {
		return nil, err
	}

	return &candidate, nil
}

func getJobApplicationRequirement(
	ctx context.Context,
	applicationID int64,
) (*model.JobApplicationRequirement, error) {
	var jobApplicationRequirement model.JobApplicationRequirement

	var itemsString []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT application_id,items,match_value,created_at,updated_at FROM job_application_requirements WHERE application_id = ?`,
		applicationID,
	).Scan(
		&jobApplicationRequirement.ApplicationID,
		&itemsString,
		&jobApplicationRequirement.MatchValue,
		&jobApplicationRequirement.CreatedAt,
		&jobApplicationRequirement.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job application requirements")

		return nil, err
	}

	err = json.Unmarshal(itemsString, &jobApplicationRequirement.Items)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal application requirements")

		return nil, err
	}

	return &jobApplicationRequirement, nil
}

func getJobApplicationCulturalFit(
	ctx context.Context,
	applicationID int64,
	jobID int64,
) (*model.JobApplicationCulturalFit, error) {
	var jobApplicationCulturalFit model.JobApplicationCulturalFit

	var answersString []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT application_id,answers,match_value,created_at,updated_at FROM job_application_cultural_fit WHERE application_id = ?`,
		applicationID,
	).Scan(
		&jobApplicationCulturalFit.ApplicationID,
		&answersString,
		&jobApplicationCulturalFit.MatchValue,
		&jobApplicationCulturalFit.CreatedAt,
		&jobApplicationCulturalFit.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job application cultural fit")

		return nil, err
	}

	err = json.Unmarshal(answersString, &jobApplicationCulturalFit.Answers)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal application cultural fit")

		return nil, err
	}

	jobApplicationCulturalFit.JobCulturalFit, err = getJobCulturalFit(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return &jobApplicationCulturalFit, nil
}

func getJobApplicationQuestion(
	ctx context.Context,
	applicationID int64,
	jobID int64,
) (*model.JobApplicationQuestion, error) {
	var jobApplicationQuestion model.JobApplicationQuestion

	var questionsString []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT application_id,questions,score,open_field_score,created_at,updated_at FROM job_application_questions WHERE application_id = ?`,
		applicationID,
	).Scan(
		&jobApplicationQuestion.ApplicationID,
		&questionsString,
		&jobApplicationQuestion.Score,
		&jobApplicationQuestion.OpenFieldScore,
		&jobApplicationQuestion.CreatedAt,
		&jobApplicationQuestion.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job application question")

		return nil, err
	}

	err = json.Unmarshal(questionsString, &jobApplicationQuestion.Questions)
	if err != nil {
		logrus.WithError(err).Error("Error to unmarshal application question")

		return nil, err
	}

	questionsMap, err := questions_adp.ListQuestions(ctx, jobID)
	if err != nil {
		return nil, err
	}

	jobApplicationQuestion.JobQuestions = make([]model.Question, 0)

	for _, question := range questionsMap {
		jobApplicationQuestion.JobQuestions = append(jobApplicationQuestion.JobQuestions, question)
	}

	return &jobApplicationQuestion, nil
}

func getJobApplicationCandidateVideo(
	ctx context.Context,
	applicationID int64,
) (*model.JobApplicationCandidateVideo, error) {
	var jobApplicationCandidateVideo model.JobApplicationCandidateVideo

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT application_id,score,created_at,updated_at FROM job_application_candidate_videos WHERE application_id = ?`,
		applicationID,
	).Scan(
		&jobApplicationCandidateVideo.ApplicationID,
		&jobApplicationCandidateVideo.Score,
		&jobApplicationCandidateVideo.CreatedAt,
		&jobApplicationCandidateVideo.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		logrus.WithError(err).Error("Error to get job application candidate video")

		return nil, err
	}

	return &jobApplicationCandidateVideo, nil
}

func getJobCulturalFit(
	ctx context.Context,
	jobID int64,
) (*model.JobCulturalFit, error) {
	jobCulturalFit := model.JobCulturalFit{}

	var answersJSON []byte

	err := database.Database.QueryRowContext(
		ctx,
		`SELECT answers 
				FROM job_cultural_fit WHERE job_id = ?`,
		jobID,
	).Scan(&answersJSON)
	if err != nil {
		logrus.WithError(err).Error("Error to get job cultural answers")

		return nil, err
	}

	if err = json.Unmarshal(answersJSON, &jobCulturalFit.Answers); err != nil {
		logrus.WithError(err).Error("Error to unmarshal job cultural answers")

		return nil, err
	}

	return &jobCulturalFit, nil
}
