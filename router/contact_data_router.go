package router

import (
	"github.com/AutoVentures-New/GO-API/handler/contact-data"
	"github.com/gofiber/fiber/v2"
)

func setupContactDataRoute(router fiber.Router) {
	router.Get("/list-data", contact_data.ListContactData)
}
