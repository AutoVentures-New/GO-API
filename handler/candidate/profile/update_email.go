package profile

import (
	"github.com/gofiber/fiber/v2"
	candidate_auth_adp "github.com/hubjob/api/app/adapters/candidate/auth"
	"github.com/hubjob/api/app/adapters/candidate/profile"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type UpdateCandidateEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func UpdateCandidateEmail(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)
	request := new(UpdateCandidateEmailRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	isValidCode, err := candidate_auth_adp.IsValidCode(
		fiberCtx.UserContext(),
		request.Code,
		request.Email,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	if !isValidCode {
		return responses.NotFound(fiberCtx, "Invalid code")
	}

	candidate.Email = request.Email

	candidate, err = profile.UpdateEmail(fiberCtx.UserContext(), candidate, request.Code)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, candidate)
}
