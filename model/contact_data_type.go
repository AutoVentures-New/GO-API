package model

type ContactDataType string

const (
	TypeEmail           ContactDataType = "EMAIL"
	TypeNote            ContactDataType = "NOTE"
	TypeFile            ContactDataType = "FILE"
	TypeSMS             ContactDataType = "SMS"
	TypeCall            ContactDataType = "CALL"
	TypeEvent           ContactDataType = "EVENT"
	TypeTask            ContactDataType = "TASK"
	TypeMeeting         ContactDataType = "MEETING"
	TypeVideoConference ContactDataType = "VIDEO_CONFERENCE"
	TypePhoneCall       ContactDataType = "PHONE_CALL"
)
