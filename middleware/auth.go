package middleware

import (
	"encoding/json"
	"fmt"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/handler/responses"
	"github.com/hubjob/api/model"
	"github.com/hubjob/api/pkg"
)

func ProtectedCompany() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config.JwtSecret)},
		ErrorHandler: func(fiberCtx *fiber.Ctx, err error) error {
			return responses.Forbidden(fiberCtx)
		},
		SuccessHandler: func(fiberCtx *fiber.Ctx) error {
			token := fiberCtx.Locals("user").(*jwt.Token)

			value, err := pkg.SessionClient.Get(fiberCtx.UserContext(), pkg.FormatToken(token.Raw)).Result()
			if err != nil {
				return responses.Forbidden(fiberCtx)
			}

			user := model.User{}
			_ = json.Unmarshal([]byte(value), &user)

			fiberCtx.Locals("user", user)

			return fiberCtx.Next()
		},
	})
}

func ProtectedCandidate(ignoreError bool) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config.JwtSecret)},
		ErrorHandler: func(fiberCtx *fiber.Ctx, err error) error {
			if ignoreError {
				return fiberCtx.Next()
			}

			return responses.Forbidden(fiberCtx)
		},
		SuccessHandler: func(fiberCtx *fiber.Ctx) error {
			token := fiberCtx.Locals("user").(*jwt.Token)

			value, err := pkg.SessionClient.Get(fiberCtx.UserContext(), pkg.FormatToken(token.Raw)).Result()
			if err != nil && !ignoreError {
				fmt.Println("entrou aqui", ignoreError)
				return responses.Forbidden(fiberCtx)
			}

			candidate := model.Candidate{}
			_ = json.Unmarshal([]byte(value), &candidate)

			fiberCtx.Locals("candidate", candidate)

			return fiberCtx.Next()
		},
	})
}
