package company

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	company_job_adp "github.com/hubjob/api/app/adapters/company/job"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type CreateJobRequest struct {
	Title               string    `json:"title"`
	IsTalentBank        bool      `json:"is_talent_bank"`
	IsSpecialNeeds      bool      `json:"is_special_needs"`
	Description         string    `json:"description"`
	JobMode             string    `json:"job_mode"`
	ContractingModality string    `json:"contracting_modality"`
	State               string    `json:"state"`
	City                string    `json:"city"`
	Responsibilities    string    `json:"responsibilities"`
	Questionnaire       string    `json:"questionnaire"`
	VideoLink           string    `json:"video_link"`
	Status              string    `json:"status"`
	PublishAt           time.Time `json:"publish_at"`
	FinishAt            time.Time `json:"finish_at"`
	JobCulturalFit      struct {
		Answers []struct {
			CulturalFitID int64  `json:"cultural_fit_id"`
			Answer        string `json:"answer"`
		} `json:"answers"`
	} `json:"cultural_fit"`
	JobRequirements struct {
		MinMatch int64 `json:"min_match"`
		Items    []struct {
			Name     string `json:"name"`
			Required bool   `json:"required"`
		} `json:"items"`
	} `json:"requirements"`
	Benefits       []int64 `json:"benefits"`
	VideoQuestions struct {
		Questions []string `json:"questions"`
	} `json:"video_questions"`
	Questions []int64 `json:"questions"`
}

func CreateJob(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	request := new(CreateJobRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	jobCulturalFit := model.JobCulturalFit{}
	jobCulturalFit.Answers = make([]model.JobCulturalFitAnswer, 0)
	for _, job := range request.JobCulturalFit.Answers {
		jobCulturalFit.Answers = append(jobCulturalFit.Answers, model.JobCulturalFitAnswer{
			CulturalFitID: job.CulturalFitID,
			Answer:        job.Answer,
		})
	}

	jobRequirement := model.JobRequirement{
		MinMatch: request.JobRequirements.MinMatch,
	}
	jobRequirement.Items = make([]model.JobRequirementItem, 0)
	for index, value := range request.JobRequirements.Items {
		jobRequirement.Items = append(jobRequirement.Items, model.JobRequirementItem{
			ID:       int64(index + 1),
			Name:     value.Name,
			Required: value.Required,
		})
	}

	benefits := make([]model.Benefit, 0)
	for _, benefit := range request.Benefits {
		benefits = append(benefits, model.Benefit{
			ID: benefit,
		})
	}

	questions := make([]model.Question, 0)
	for _, question := range request.Questions {
		questions = append(questions, model.Question{
			ID: question,
		})
	}

	job, err := company_job_adp.CreateJob(
		fiberCtx.UserContext(),
		model.Job{
			Title:               request.Title,
			CompanyID:           user.CompanyID,
			IsTalentBank:        request.IsTalentBank,
			IsSpecialNeeds:      request.IsSpecialNeeds,
			Description:         request.Description,
			JobMode:             request.JobMode,
			ContractingModality: request.ContractingModality,
			State:               request.State,
			City:                request.City,
			Responsibilities:    request.Responsibilities,
			Questionnaire:       request.Questionnaire,
			VideoLink:           request.VideoLink,
			Status:              request.Status,
			PublishAt:           request.PublishAt,
			FinishAt:            request.FinishAt,
			JobCulturalFit:      &jobCulturalFit,
			JobRequirement:      &jobRequirement,
			Benefits:            benefits,
			VideoQuestions:      &model.JobVideoQuestions{Questions: request.VideoQuestions.Questions},
			Questions:           questions,
		},
	)
	if errors.Is(err, company_job_adp.ErrJobAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY|JOB|ALREADY_EXISTS", "Job already exists")
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, job)
}

func ListJobs(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	jobs, err := company_job_adp.ListJobs(
		fiberCtx.UserContext(),
		user.CompanyID,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, jobs)
}

func GetJob(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	job, err := company_job_adp.GetJob(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_job_adp.ErrJobNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, job)
}

type UpdateJobRequest struct {
	Title               string    `json:"title"`
	IsTalentBank        bool      `json:"is_talent_bank"`
	IsSpecialNeeds      bool      `json:"is_special_needs"`
	Description         string    `json:"description"`
	JobMode             string    `json:"job_mode"`
	ContractingModality string    `json:"contracting_modality"`
	State               string    `json:"state"`
	City                string    `json:"city"`
	Responsibilities    string    `json:"responsibilities"`
	Questionnaire       string    `json:"questionnaire"`
	VideoLink           string    `json:"video_link"`
	Status              string    `json:"status"`
	PublishAt           time.Time `json:"publish_at"`
	FinishAt            time.Time `json:"finish_at"`
	JobCulturalFit      struct {
		Answers []struct {
			CulturalFitID int64  `json:"cultural_fit_id"`
			Answer        string `json:"answer"`
		} `json:"answers"`
	} `json:"cultural_fit"`
	JobRequirements struct {
		MinMatch int64 `json:"min_match"`
		Items    []struct {
			Name     string `json:"name"`
			Required bool   `json:"required"`
		} `json:"items"`
	} `json:"requirements"`
	Benefits       []int64 `json:"benefits"`
	VideoQuestions struct {
		Questions []string `json:"questions"`
	} `json:"video_questions"`
	Questions []int64 `json:"questions"`
}

func UpdateJob(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	request := new(UpdateJobRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	job, err := company_job_adp.GetJob(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_job_adp.ErrJobNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	job.JobCulturalFit.Answers = make([]model.JobCulturalFitAnswer, 0)
	for _, value := range request.JobCulturalFit.Answers {
		job.JobCulturalFit.Answers = append(job.JobCulturalFit.Answers, model.JobCulturalFitAnswer{
			CulturalFitID: value.CulturalFitID,
			Answer:        value.Answer,
		})
	}

	job.JobRequirement.MinMatch = request.JobRequirements.MinMatch

	job.JobRequirement.Items = make([]model.JobRequirementItem, 0)
	for index, value := range request.JobRequirements.Items {
		job.JobRequirement.Items = append(job.JobRequirement.Items, model.JobRequirementItem{
			ID:       int64(index + 1),
			Name:     value.Name,
			Required: value.Required,
		})
	}

	benefits := make([]model.Benefit, 0)
	for _, benefit := range request.Benefits {
		benefits = append(benefits, model.Benefit{
			ID: benefit,
		})
	}

	questions := make([]model.Question, 0)
	for _, question := range request.Questions {
		questions = append(questions, model.Question{
			ID: question,
		})
	}

	job.Title = request.Title
	job.IsTalentBank = request.IsTalentBank
	job.IsSpecialNeeds = request.IsSpecialNeeds
	job.Description = request.Description
	job.JobMode = request.JobMode
	job.ContractingModality = request.ContractingModality
	job.State = request.State
	job.City = request.City
	job.Responsibilities = request.Responsibilities
	job.Questionnaire = request.Questionnaire
	job.VideoLink = request.VideoLink
	job.Status = request.Status
	job.PublishAt = request.PublishAt
	job.FinishAt = request.FinishAt
	job.Benefits = benefits
	job.VideoQuestions.Questions = request.VideoQuestions.Questions
	job.Questions = questions

	job, err = company_job_adp.UpdateJob(
		fiberCtx.UserContext(),
		job,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, job)
}

func DeleteJob(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	err = company_job_adp.DeleteJob(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_job_adp.ErrJobNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}
