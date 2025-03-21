package company

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	company_user_adp "github.com/hubjob/api/app/adapters/company/user"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
	"strconv"
)

func ListUsers(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	users, err := company_user_adp.ListUsers(
		fiberCtx.UserContext(),
		user.CompanyID,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, users)
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
}

func CreateUser(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	request := new(CreateUserRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	userModel, err := company_user_adp.CreateUser(
		fiberCtx.UserContext(),
		model.User{
			Name:      request.Name,
			CompanyID: user.CompanyID,
			CPF:       request.CPF,
			Email:     request.Email,
		},
	)
	if errors.Is(err, company_user_adp.ErrUserEmailAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY|USERS|EMAIL_ALREADY_EXISTS", "User email already exists")
	}

	if errors.Is(err, company_user_adp.ErrUserCPFAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY|USERS|CPF_ALREADY_EXISTS", "User cpf already exists")
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, userModel)
}

func DeleteUser(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	id := fiberCtx.Params("id")
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Params {id} is required")
	}

	idInt, err := strconv.Atoi(id)
	if len(id) == 0 {
		return responses.BadRequest(fiberCtx, "Invalid params {id}")
	}

	err = company_user_adp.DeleteUser(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_user_adp.ErrUserNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}

type UpdateUserRequest struct {
	Name   string       `json:"name"`
	CPF    string       `json:"cpf"`
	Email  string       `json:"email"`
	Status model.Status `json:"status"`
}

func UpdateUser(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	request := new(UpdateUserRequest)

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

	userModel, err := company_user_adp.GetUser(
		fiberCtx.UserContext(),
		int64(idInt),
		user.CompanyID,
	)
	if errors.Is(err, company_user_adp.ErrUserNotFound) {
		return responses.NotFound(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	userModel.Name = request.Name
	userModel.Email = request.Email
	userModel.CPF = request.CPF
	userModel.Status = request.Status
	userModel, err = company_user_adp.UpdateUser(
		fiberCtx.UserContext(),
		userModel,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, userModel)
}
