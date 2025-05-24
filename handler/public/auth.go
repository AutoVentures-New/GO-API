package public

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hubjob/api/app/adapters/company/auth"
	"github.com/hubjob/api/app/adapters/public"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
	"github.com/hubjob/api/pkg"
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
		return responses.NotFound(fiberCtx, err.Error())
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

type CreatePasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func CreatePassword(fiberCtx *fiber.Ctx) error {
	request := CreatePasswordRequest{}

	err := fiberCtx.BodyParser(&request)
	if err != nil {
		return responses.BadRequest(fiberCtx, "Invalid request")
	}

	if request.Token == "" || request.Password == "" {
		return responses.BadRequest(fiberCtx, "Invalid request")
	}

	email, err := public.CreatePassword(
		fiberCtx.UserContext(),
		request.Token,
		request.Password,
	)
	if errors.Is(err, public.ErrUserTokenNotFound) {
		return responses.BadRequest(fiberCtx, err.Error())
	}

	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	user, err := auth.GetUser(fiberCtx.UserContext(), email)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
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
