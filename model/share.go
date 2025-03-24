package model

type Status string

const (
	ACTIVE   Status = "ACTIVE"
	INACTIVE Status = "INACTIVE"
	PENDING  Status = "PENDING"
)

const NULL_PASSWORD = "NULL_PASSWORD"
