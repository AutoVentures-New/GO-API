package public

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/app/adapters/public"
	"github.com/hubjob/api/handler/responses"
)

type ForgotPasswordRequest struct {
	Email    string `json:"email"`
	ExecType string `json:"exec_type"`
}

func ForgotPassword(fiberCtx *fiber.Ctx) error {
	request := ForgotPasswordRequest{}

	err := fiberCtx.BodyParser(&request)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid request")
	}

	if request.ExecType != "COMPANY" && request.ExecType != "CANDIDATE" {
		return responses.BadRequest(fiberCtx, "Invalid exec_type")
	}

	if request.Email == "" {
		return responses.BadRequest(fiberCtx, "Email is required")
	}

	err = public.ForgotPassword(
		fiberCtx.UserContext(),
		request.Email,
		request.ExecType,
	)
	if errors.Is(err, public.ErrEmailNotFound) {
		return responses.BadRequest(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}

type ChangePasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func ChangePassword(fiberCtx *fiber.Ctx) error {
	request := ChangePasswordRequest{}

	err := fiberCtx.BodyParser(&request)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid request")
	}

	if request.Token == "" || request.Password == "" {
		return responses.BadRequest(fiberCtx, "Invalid request")
	}

	err = public.ChangePassword(
		fiberCtx.UserContext(),
		request.Token,
		request.Password,
	)
	if errors.Is(err, public.ErrTokenNotFound) {
		return responses.BadRequest(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}
