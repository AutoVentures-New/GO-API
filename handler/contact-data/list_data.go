package contact_data

import (
	"github.com/AutoVentures-New/GO-API/app/adapters/contact-data"
	"github.com/AutoVentures-New/GO-API/handler/responses"
	"github.com/AutoVentures-New/GO-API/model"
	"github.com/AutoVentures-New/GO-API/model/request"
	"github.com/AutoVentures-New/GO-API/pkg"
	"github.com/gofiber/fiber/v2"
	"sync"
)

type ListResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
}

func ListSuccess[T any](ctx *fiber.Ctx, data []T, total int) error {
	return ctx.Status(fiber.StatusOK).JSON(ListResponse[T]{
		Data:  data,
		Total: total,
	})
}

func ListContactData(fiberCtx *fiber.Ctx) error {
	user := fiberCtx.Locals("user").(model.User)
	ctx := fiberCtx.Context()
	account := user.Account
	var query request.ContactDataQuery

	if err := fiberCtx.QueryParser(&query); err != nil {
		return responses.BadRequest(fiberCtx, "Invalid query parameters")
	}

	if query.Type == "" {
		query.Type = model.TypeAll
	}

	contactData, err := contact_data.GetContactData(ctx, account, query)
	if err != nil {
		return responses.InternalServerError(fiberCtx, err)
	}

	if len(contactData) == 0 {
		return ListSuccess(fiberCtx, contactData, 0)
	}

	ulids := pkg.ExtractIdentifiers(contactData)

	var (
		emails        []model.EmailBucket
		calls         []model.Call
		activityFiles []model.ActivityFile
		notes         []model.Note
		events        []model.CalendarEvent
	)

	var (
		wg      sync.WaitGroup
		errLock sync.Mutex
		callErr error
	)

	fetch := func(do bool, action func() error) {
		if do {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := action(); err != nil {
					errLock.Lock()
					defer errLock.Unlock()
					if callErr == nil {
						callErr = err
					}
				}
			}()
		}
	}

	fetch(query.Type == model.TypeAll || query.Type == model.TypeEmail, func() error {
		var err error
		emails, err = contact_data.GetEmails(ctx, account, ulids)
		return err
	})

	fetch(query.Type == model.TypeAll || query.Type == model.TypeCall || query.Type == model.TypePhoneCall, func() error {
		var err error
		calls, err = contact_data.GetCalls(ctx, account, ulids, *query.ContactULID)
		return err
	})

	fetch(query.Type == model.TypeAll || query.Type == model.TypeFile, func() error {
		var err error
		activityFiles, err = contact_data.GetActivityFiles(ctx, account, ulids)
		return err
	})

	fetch(query.Type == model.TypeAll || query.Type == model.TypeNote, func() error {
		var err error
		notes, err = contact_data.GetNotes(ctx, account, ulids)
		return err
	})

	fetch(query.Type == model.TypeAll || query.Type == model.TypeCall || query.Type == model.TypeEvent || query.Type == model.TypeTask || query.Type == model.TypeMeeting || query.Type == model.TypeVideoConference, func() error {
		var err error
		events, err = contact_data.GetEvents(ctx, account, ulids)
		return err
	})

	wg.Wait()

	if callErr != nil {
		return responses.InternalServerError(fiberCtx, callErr)
	}

	contactDataMap := pkg.SliceToMap(contactData, func(c model.ContactData) string { return c.Identifier })
	for i, e := range events {
		events[i].From = contactDataMap[e.Ulid].From
		events[i].Date = contactDataMap[e.Ulid].Date
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

	return ListSuccess(fiberCtx, contactData, len(contactData))
}
