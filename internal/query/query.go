package query

var ListContactData = `SELECT cd.ulid, cd.type, cd.identifier, cd.from, cd.to, cd.cc, cd.date FROM tenant_%s.contact_data_contact cdc INNER JOIN tenant_%s.contact_data cd ON cdc.contact_data_ulid COLLATE utf8mb4_unicode_ci = cd.ulid COLLATE utf8mb4_unicode_ci WHERE 1 = 1 `

var ListEmailData = `SELECT record_number, ulid, account_id, message_id, thread_id, subject, eb.from, eb.to, cc, bcc, reply_to, headers, starred, unread, reply_to_message_id, body, files, folder, links, opens, COALESCE(link_clicks, 'null') AS link_clicks, is_tracked, eb.date, created_at, updated_at FROM tenant_%s.email_bucket eb`

var ListCallData = `SELECT record_number, ulid, 'Phone Call' as subject, created_by, user_phone_number, contact_phone_number, done, c.to, call_type, direction, outcome, notes, user_create_date, created_at, updated_at FROM tenant_%s.calls c`

var ListActivityFileData = `SELECT record_number, ulid, created_by, af.to, subject, done, files, user_create_date, created_at, updated_at FROM tenant_%s.activities_files af`

var ListNoteData = `SELECT record_number, ulid, created_by, n.to, subject, done, text, commented_at, files, user_create_date, created_at, updated_at FROM tenant_%s.notes n`

var ListCalendarEventData = `SELECT record_number, ulid, name, description, participants, ce.when, location, recurrence, notifications, conferencing, conference_records as 'records', organizer_name, organizer_email, owner, done, start_date, end_date, all_day, type, sequence, files, created_at, updated_at FROM tenant_%s.calendar_events ce`

var ListCommentsData = `SELECT record_number, ulid, created_by, text, commented_at, files, created_at, updated_at FROM tenant_%s.comments`

var ListEventFileData = `SELECT record_number, ulid, event_ulid, is_external, name, extension, link, created_at, updated_at FROM tenant_%s.calendar_event_files`
