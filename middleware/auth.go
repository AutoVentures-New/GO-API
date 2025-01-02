package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/config"
	"github.com/hubjob/api/handler/responses"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config.JwtSecret)},
		ErrorHandler: func(fiberCtx *fiber.Ctx, err error) error {
			return responses.Forbidden(fiberCtx)
		},
	})
}
