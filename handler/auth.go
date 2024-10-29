package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/trabalhe-conosco/api/config"
	"github.com/trabalhe-conosco/api/database"
	"github.com/trabalhe-conosco/api/model"
	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	return string(bytes), err
}

func getUserByEmail(ctx context.Context, email string) (model.User, error) {
	user := model.User{}

	err := database.Database.QueryRowContext(
		ctx,
		"SELECT * FROM b2b_users where email = ?",
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		logrus.WithError(err).
			WithField("email", email).
			Error("Error to find user by email")

		return model.User{}, err
	}

	return user, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	request := new(LoginRequest)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	user, err := getUserByEmail(c.UserContext(), request.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	}

	if !checkPasswordHash(request.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or passwordaaa", "data": nil})
	}

	fmt.Println(user)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().UTC().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config.JwtSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(map[string]interface{}{
		"data": map[string]interface{}{
			"token":       t,
			"expire_date": claims["exp"],
		},
	})
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(fiberCtx *fiber.Ctx) error {
	request := new(RegisterRequest)

	if err := fiberCtx.BodyParser(&request); err != nil {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request"})
	}

	ctx := fiberCtx.UserContext()
	now := time.Now().UTC()

	password, err := hashPassword(request.Password)
	if err != nil {
		logrus.WithError(err).Error("Error on create password")

		return fiberCtx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error"})
	}

	var mysqlErr *mysql.MySQLError

	_, err = database.Database.ExecContext(
		ctx,
		"INSERT INTO b2b_users VALUES(0, ?, ?, ?, 'ACTIVE', ?, ?)",
		request.Name,
		request.Email,
		password,
		now,
		now,
	)
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return fiberCtx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Email already exists"})
	}

	if err != nil {
		logrus.WithError(err).Error("Error on insert user")

		return fiberCtx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = request.Email
	claims["exp"] = time.Now().UTC().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config.JwtSecret))
	if err != nil {
		return fiberCtx.SendStatus(fiber.StatusInternalServerError)
	}

	return fiberCtx.JSON(map[string]interface{}{
		"data": map[string]interface{}{
			"token":       t,
			"expire_date": claims["exp"],
		},
	})
}
