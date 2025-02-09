package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hubjob/api/handler/configuration"
)

func setupConfigurationRoute(router fiber.Router) {
	router.Get("/cultural-fit", configuration.GetCulturalFit)
	router.Get("/questionnaires/behavioral", configuration.GetQuestionnaireBehavioral)
	router.Get("/questionnaires/professional", configuration.GetQuestionnaireProfessional)
	router.Get("/cultural-fit", configuration.GetCulturalFit)
	router.Get("/list-areas", configuration.ListAreas)
}
