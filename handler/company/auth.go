package company

import (
	"encoding/json"
	"errors"
	"github.com/hubjob/api/pkg"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	company_auth_adp "github.com/hubjob/api/app/adapters/company/auth"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
	"golang.org/x/crypto/bcrypt"
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
	if errors.Is(err, company_auth_adp.ErrUserAlreadyExists) {
		return responses.Conflict(fiberCtx, "USER_ALREADY_EXISTS", err.Error())
	}

	if errors.Is(err, company_auth_adp.ErrCompanyAlreadyExists) {
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

	return createToken(fiberCtx, userResponse)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(fiberCtx *fiber.Ctx) error {
	request := new(LoginRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	user, err := company_auth_adp.GetUser(fiberCtx.UserContext(), request.Email)
	if err != nil {
		return responses.Unauthorized(fiberCtx)
	}

	if !checkPasswordHash(request.Password, user.Password) {
		return responses.Unauthorized(fiberCtx)
	}

	return createToken(fiberCtx, user)
}

func createToken(fiberCtx *fiber.Ctx, user model.User) error {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().UTC().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config.JwtSecret))
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	userString, _ := json.Marshal(user)

	err = pkg.SessionClient.Set(fiberCtx.UserContext(), pkg.FormatToken(t), string(userString), time.Hour*72).Err()
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, map[string]interface{}{
		"token":       t,
		"expire_date": claims["exp"],
	})
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func Me(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	user, err := company_auth_adp.GetUser(fiberCtx.UserContext(), user.Email)
	if err != nil {
		return responses.Forbidden(fiberCtx)
	}

	return responses.Success(fiberCtx, user)
}
