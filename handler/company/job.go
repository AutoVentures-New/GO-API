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
}

func CreateJob(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	request := new(CreateJobRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
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

	job, err := company_job_adp.UpdateJob(
		fiberCtx.UserContext(),
		model.Job{
			ID:                  int64(idInt),
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
		},
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
