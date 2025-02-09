package candidate

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	candidate_job_adp "github.com/hubjob/api/app/adapters/candidate/job"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func StartApplication(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	id := fiberCtx.Params("job_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {job_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {job_id}")
	}

	application, err := candidate_job_adp.StartApplication(
		fiberCtx.UserContext(),
		candidate.ID,
		int64(idInt),
	)
	if errors.Is(err, candidate_job_adp.ErrJobNotFound) {
		return responses.NotFound(fiberCtx, "Job not found")
	}

	if errors.Is(err, candidate_job_adp.ErrInvalidJob) {
		return responses.BadRequest(fiberCtx, err.Error())
	}

	if errors.Is(err, candidate_job_adp.ErrApplicationAlreadyExist) {
		return responses.Conflict(fiberCtx, "CANDIDATE|START_APPLICATION|ALREADY_EXIST", "Application already exist")
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, application)
}

func GetApplication(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	id := fiberCtx.Params("job_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {job_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {job_id}")
	}

	application, err := candidate_job_adp.GetJobApplication(
		fiberCtx.UserContext(),
		int64(idInt),
		candidate.ID,
	)
	if errors.Is(err, candidate_job_adp.ErrApplicationNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, application)
}

func CanceledApplication(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	id := fiberCtx.Params("job_id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {job_id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid params {job_id}")
	}

	application, err := candidate_job_adp.GetJobApplication(
		fiberCtx.UserContext(),
		int64(idInt),
		candidate.ID,
	)
	if errors.Is(err, candidate_job_adp.ErrApplicationNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	if application.Status != model.FILLING && application.Status != model.WAITING_EVALUATION {
		return responses.Success(fiberCtx, application)
	}

	application, err = candidate_job_adp.CanceledApplication(fiberCtx.UserContext(), application)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, application)
}

func ListApplications(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	applications, err := candidate_job_adp.ListJobApplications(
		fiberCtx.UserContext(),
		candidate.ID,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, applications)
}
