package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/AutoVentures-New/GO-API/handler/contact-data"
)

func setupContactDataRoute(router fiber.Router) {
	router.Get("/", contact_data.ListContactData)
}
