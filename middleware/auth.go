package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/AutoVentures-New/GO-API/config"
	"github.com/AutoVentures-New/GO-API/handler/responses"
	"github.com/AutoVentures-New/GO-API/model"
)

func Auth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		accessToken := ctx.Get("access_token")

		if accessToken != config.Config.AccessToken {
			return responses.Forbidden(ctx)
		}

		user := model.User{
			User:    ctx.Get("user", ""),
			Account: ctx.Get("account", ""),
		}

		ctx.Locals("user", user)

		return ctx.Next()
	}
}
