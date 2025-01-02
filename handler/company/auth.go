package company

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	company_auth_adp "github.com/hubjob/api/app/adapters/company/auth"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

type ValidateCnpjCpfRequest struct {
	Cnpj string `json:"cnpj"`
	Cpf  string `json:"cpf"`
}

func ValidateCnpjCpf(fiberCtx *fiber.Ctx) error {
	request := new(ValidateCnpjCpfRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	err := company_auth_adp.ValidateCnpjCpf(
		fiberCtx.UserContext(),
		request.Cnpj,
		request.Cpf,
	)
	if errors.Is(err, company_auth_adp.CpfAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY|AUTH|CPF_ALREADY_EXISTS", "Cpf already exists")
	}

	if errors.Is(err, company_auth_adp.CnpjAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY|AUTH|CNPJ_ALREADY_EXISTS", "Cnpj already exists")
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}

type SendEmailValidationRequest struct {
	Email string `json:"email"`
}

func SendEmailValidation(fiberCtx *fiber.Ctx) error {
	request := new(SendEmailValidationRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	err := company_auth_adp.SendEmailValidation(
		fiberCtx.UserContext(),
		request.Email,
	)
	if errors.Is(err, company_auth_adp.EmailAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY|AUTH|EMAIL_ALREADY_EXISTS", "Email already exists")
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, nil)
}

type ValidateEmailValidationCodeRequest struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

func ValidateEmailValidationCode(fiberCtx *fiber.Ctx) error {
	request := new(ValidateEmailValidationCodeRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	isValidCode, err := company_auth_adp.IsValidCode(
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

	return responses.Success(fiberCtx, nil)
}

type CreateAccountRequest struct {
	Name     string `json:"name"`
	CNPJ     string `json:"cnpj"`
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`

	City  string `json:"city"`
	State string `json:"state"`
}

func CreateAccount(fiberCtx *fiber.Ctx) error {
	request := new(CreateAccountRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	isValidCode, err := company_auth_adp.IsValidCode(
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

	user := model.User{
		Name:     request.Name,
		CPF:      request.CPF,
		Email:    request.Email,
		Password: request.Password,
	}

	company := model.Company{
		CNPJ: request.CNPJ,
	}

	err = company_auth_adp.CheckAlreadyExist(
		fiberCtx.UserContext(),
		user,
		company,
	)
	if errors.As(err, &company_auth_adp.ErrUserAlreadyExists) {
		return responses.Conflict(fiberCtx, "USER_ALREADY_EXISTS", err.Error())
	}

	if errors.As(err, &company_auth_adp.ErrCompanyAlreadyExists) {
		return responses.Conflict(fiberCtx, "COMPANY_ALREADY_EXISTS", err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	userResponse, err := company_auth_adp.CreateAccount(
		fiberCtx.UserContext(),
		user,
		company,
		request.Code,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, userResponse)
}
