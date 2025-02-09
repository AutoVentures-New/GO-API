package candidate

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	candidate_auth_adp "github.com/hubjob/api/app/adapters/candidate/auth"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
	"golang.org/x/crypto/bcrypt"
)

type ValidateCpfRequest struct {
	Cpf string `json:"cpf"`
}

func ValidateCpf(fiberCtx *fiber.Ctx) error {
	request := new(ValidateCpfRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	err := candidate_auth_adp.ValidateCpf(
		fiberCtx.UserContext(),
		request.Cpf,
	)
	if errors.Is(err, candidate_auth_adp.CpfAlreadyExists) {
		return responses.Conflict(fiberCtx, "CANDIDATE|AUTH|CPF_ALREADY_EXISTS", "Cpf already exists")
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

	err := candidate_auth_adp.SendEmailValidation(
		fiberCtx.UserContext(),
		request.Email,
	)
	if errors.Is(err, candidate_auth_adp.EmailAlreadyExists) {
		return responses.Conflict(fiberCtx, "CANDIDATE|AUTH|EMAIL_ALREADY_EXISTS", "Email already exists")
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

	return responses.Success(fiberCtx, nil)
}

type CreateAccountRequest struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birth_date"`
	State     string `json:"state"`
	City      string `json:"city"`
}

func CreateAccount(fiberCtx *fiber.Ctx) error {
	request := new(CreateAccountRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return responses.InvalidBodyRequest(fiberCtx, err)
	}

	birthDate, err := time.Parse("2006-01-02", request.BirthDate)
	if err != nil {
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

	candidate := model.Candidate{
		Name:      request.Name,
		CPF:       request.CPF,
		Email:     request.Email,
		Phone:     request.Phone,
		BirthDate: birthDate,
		Password:  request.Password,
		Address: model.Address{
			State: request.State,
			City:  request.City,
		},
	}

	err = candidate_auth_adp.CheckAlreadyExist(
		fiberCtx.UserContext(),
		candidate,
	)
	if errors.Is(err, candidate_auth_adp.ErrCandidateAlreadyExists) {
		return responses.Conflict(fiberCtx, "CANDIDATE_ALREADY_EXISTS", err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	candidateResponse, err := candidate_auth_adp.CreateAccount(
		fiberCtx.UserContext(),
		candidate,
		request.Code,
	)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return createToken(fiberCtx, candidateResponse)
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

	candidate, err := candidate_auth_adp.GetCandidate(fiberCtx.UserContext(), request.Email)
	if err != nil {
		return responses.Unauthorized(fiberCtx)
	}

	if !checkPasswordHash(request.Password, candidate.Password) {
		return responses.Unauthorized(fiberCtx)
	}

	return createToken(fiberCtx, candidate)
}

func createToken(fiberCtx *fiber.Ctx, candidate model.Candidate) error {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["candidate"] = candidate
	claims["exp"] = time.Now().UTC().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config.JwtSecret))
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
	candidate := fiberCtx.Locals("candidate").(model.Candidate)

	candidate, err := candidate_auth_adp.GetCandidate(fiberCtx.UserContext(), candidate.Email)
	if err != nil {
		return responses.Forbidden(fiberCtx)
	}

	return responses.Success(fiberCtx, candidate)
}
