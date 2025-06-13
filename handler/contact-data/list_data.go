package contact_data

import (
	"github.com/AutoVentures-New/GO-API/app/adapters/contact-data"
	"github.com/AutoVentures-New/GO-API/handler/responses"
	"github.com/AutoVentures-New/GO-API/model"
	"github.com/AutoVentures-New/GO-API/model/request"
	"github.com/AutoVentures-New/GO-API/pkg"
	"github.com/gofiber/fiber/v2"
)

func ListContactData(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	ctx := fiberCtx.Context()
	account := user.Account
	var query request.ContactDataQuery

	if err := fiberCtx.QueryParser(&query); err != nil {
		return responses.BadRequest(fiberCtx, "Invalid query parameters")
	}

	contactData, err := contact_data.GetContactData(ctx, account, query)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	ulids := pkg.ExtractIdentifiers(contactData)

	emails, err := contact_data.GetEmails(ctx, account, ulids)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	calls, err := contact_data.GetCalls(ctx, account, ulids)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	activityFiles, err := contact_data.GetActivityFiles(ctx, account, ulids)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	notes, err := contact_data.GetNotes(ctx, account, ulids)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	events, err := contact_data.GetEvents(ctx, account, ulids)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	maps := map[model.ContactDataType]map[string]any{
		model.TypeEmail:           pkg.SliceToAnyMap(emails, func(e model.EmailBucket) string { return e.ThreadID }),
		model.TypePhoneCall:       pkg.SliceToAnyMap(calls, func(c model.Call) string { return c.Ulid }),
		model.TypeFile:            pkg.SliceToAnyMap(activityFiles, func(a model.ActivityFile) string { return a.Ulid }),
		model.TypeNote:            pkg.SliceToAnyMap(notes, func(n model.Note) string { return n.Ulid }),
		model.TypeCall:            pkg.SliceToAnyMap(events, func(e model.CalendarEvent) string { return e.Ulid }),
		model.TypeEvent:           pkg.SliceToAnyMap(events, func(e model.CalendarEvent) string { return e.Ulid }),
		model.TypeTask:            pkg.SliceToAnyMap(events, func(e model.CalendarEvent) string { return e.Ulid }),
		model.TypeMeeting:         pkg.SliceToAnyMap(events, func(e model.CalendarEvent) string { return e.Ulid }),
		model.TypeVideoConference: pkg.SliceToAnyMap(events, func(e model.CalendarEvent) string { return e.Ulid }),
	}

	for i, c := range contactData {
		if objMap, ok := maps[c.Type]; ok {
			if obj, exists := objMap[c.Identifier]; exists {
				contactData[i].Object = obj
			}
		}
	}

	return responses.Success(fiberCtx, contactData)
}
