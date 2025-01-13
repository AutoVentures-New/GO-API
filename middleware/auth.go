package middleware

import (
	"encoding/json"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
)

func ProtectedCompany() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config.JwtSecret)},
		ErrorHandler: func(fiberCtx *fiber.Ctx, err error) error {
			return responses.Forbidden(fiberCtx)
		},
		SuccessHandler: func(fiberCtx *fiber.Ctx) error {
			token := fiberCtx.Locals("user").(*jwt.Token)
			user := model.User{}
			userString, _ := json.Marshal(token.Claims.(jwt.MapClaims)["user"])
			_ = json.Unmarshal(userString, &user)

			fiberCtx.Locals("user", user)

			return fiberCtx.Next()
		},
	})
}

func ProtectedCandidate() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config.JwtSecret)},
		ErrorHandler: func(fiberCtx *fiber.Ctx, err error) error {
			return responses.Forbidden(fiberCtx)
		},
		SuccessHandler: func(fiberCtx *fiber.Ctx) error {
			token := fiberCtx.Locals("user").(*jwt.Token)
			candidate := model.Candidate{}
			candidateString, _ := json.Marshal(token.Claims.(jwt.MapClaims)["candidate"])
			_ = json.Unmarshal(candidateString, &candidate)

			fiberCtx.Locals("candidate", candidate)

			return fiberCtx.Next()
		},
	})
}
