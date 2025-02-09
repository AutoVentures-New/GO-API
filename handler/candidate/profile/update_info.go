package profile

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/candidate/profile"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type UpdateCandidateRequest struct {
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birth_date"`
	State     string `json:"state"`
	City      string `json:"city"`
}

func UpdateCandidate(fiberCtx *fiber.Ctx) error {
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	request := new(UpdateCandidateRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	birthDate, err := time.Parse("2006-01-02", request.BirthDate)
	if err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	candidate.Name = request.Name
	candidate.CPF = request.CPF
	candidate.Phone = request.Phone
	candidate.BirthDate = birthDate
	candidate.Address.City = request.City
	candidate.Address.State = request.State

	candidate, err = profile.UpdateInfo(fiberCtx.UserContext(), candidate)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, candidate)
}
