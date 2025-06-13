package middleware

import (
	"fmt"
	"github.com/AutoVentures-New/GO-API/config"
	"github.com/AutoVentures-New/GO-API/handler/responses"
	"github.com/AutoVentures-New/GO-API/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return responses.Forbidden(c)
		}

		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			return responses.Forbidden(c)
		}
		tokenString := authHeader[len(prefix):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(config.Config.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			return responses.Forbidden(c)
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return responses.Forbidden(c)
		}

		user := model.User{
			User:    c.Get("user", ""),
			Account: c.Get("account", ""),
		}

		c.Locals("user", user)

		return c.Next()
	}
}
