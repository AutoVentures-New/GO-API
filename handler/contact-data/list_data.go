package contact_data

import (
	"github.com/AutoVentures-New/GO-API/app/adapters/contact-data"
	"github.com/AutoVentures-New/GO-API/handler/responses"
	"github.com/AutoVentures-New/GO-API/model"
	"github.com/AutoVentures-New/GO-API/model/request"
	"github.com/gofiber/fiber/v2"
)

func ListContactData(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)

	var query request.ContactDataQuery

	if err := fiberCtx.QueryParser(&query); err != nil {
		return responses.BadRequest(fiberCtx, "Invalid query parameters")
	}

	contactData, err := contact_data.GetContactData(fiberCtx.Context(), user, query)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	_, err = contact_data.GetEmails(fiberCtx.Context(), user, ExtractIdentifiers(contactData))
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	_, err = contact_data.GetCalls(fiberCtx.Context(), user, ExtractIdentifiers(contactData))
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	_, err = contact_data.GetActivityFiles(fiberCtx.Context(), user, ExtractIdentifiers(contactData))
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	_, err = contact_data.GetNotes(fiberCtx.Context(), user, ExtractIdentifiers(contactData))
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	_, err = contact_data.GetEvents(fiberCtx.Context(), user, ExtractIdentifiers(contactData))
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	return responses.Success(fiberCtx, contactData)
}

func ExtractIdentifiers(data []model.ContactData) []string {
	identifiers := make([]string, 0, len(data))
	for _, d := range data {
		identifiers = append(identifiers, d.Identifier)
	}
	return identifiers
}
