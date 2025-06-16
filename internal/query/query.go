package query

var ListContactData = `SELECT cd.ulid, cd.type, cd.identifier, cd.from, cd.to, cd.cc, cd.date FROM tenant_%s.contact_data_contact cdc INNER JOIN tenant_%s.contact_data cd ON cdc.contact_data_ulid COLLATE utf8mb4_unicode_ci = cd.ulid COLLATE utf8mb4_unicode_ci WHERE 1 = 1 `

var CountContactData = `SELECT count(cd.ulid) FROM tenant_%s.contact_data_contact cdc INNER JOIN tenant_%s.contact_data cd ON cdc.contact_data_ulid COLLATE utf8mb4_unicode_ci = cd.ulid COLLATE utf8mb4_unicode_ci WHERE 1 = 1 `

var ListEmailData = `SELECT record_number, ulid, account_id, message_id, thread_id, subject, eb.from, eb.to, cc, bcc, reply_to, headers, starred, unread, reply_to_message_id, body, files, folder, links, opens, COALESCE(link_clicks, 'null') AS link_clicks, is_tracked, eb.date, created_at, updated_at FROM tenant_%s.email_bucket eb`

var ListCallData = `SELECT record_number, ulid, 'Phone Call' as subject, created_by, user_phone_number, contact_phone_number, done, c.to, call_type, direction, outcome, notes, user_create_date, created_at, updated_at FROM tenant_%s.calls c`

var ListActivityFileData = `SELECT record_number, ulid, created_by, af.to, subject, done, files, user_create_date, created_at, updated_at FROM tenant_%s.activities_files af`

var ListNoteData = `SELECT record_number, ulid, created_by, n.to, subject, done, text, commented_at, files, user_create_date, created_at, updated_at FROM tenant_%s.notes n`

var ListCalendarEventData = `SELECT ce.record_number, ce.ulid, cep.calendar_ulid, name, description, participants, ce.when, location, recurrence, notifications, conferencing, conference_records as 'records', organizer_name, organizer_email, owner, done, start_date, end_date, all_day, type, sequence, files, ce.created_at, ce.updated_at FROM tenant_%s.calendar_events ce INNER JOIN tenant_%s.calendar_event_providers cep ON cep.event_ulid = ce.ulid `

var ListCommentsData = `SELECT record_number, ulid, created_by, text, commented_at, files, created_at, updated_at FROM tenant_%s.comments`

var ListEventFileData = `SELECT record_number, ulid, event_ulid, is_external, name, extension, link, created_at, updated_at FROM tenant_%s.calendar_event_files`

var ListCalendarData = `SELECT ulid, user_ulid, provider_ulid, external_id, name, color FROM tenant_%s.calendars WHERE 1 = 1`

var ListUsersData = `SELECT ulid, first_name, last_name, image FROM users WHERE 1 = 1`

var ListContactsData = `SELECT contact.ulid, contact.first_name, contact.last_name, contact.company_name, COALESCE(contact.image, '') AS image FROM (SELECT p.ulid, p.first_name, p.last_name, NULL AS company_name, p.image FROM tenant_%s.people p UNION ALL SELECT c.ulid, NULL AS first_name, NULL AS last_name, c.company_name, c.image FROM tenant_%s.companies c) AS contact WHERE 1 = 1`
