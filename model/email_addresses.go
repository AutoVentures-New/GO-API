package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type EmailAddress struct {
	Email  string `json:"email"`
	Name   string `json:"name,omitempty"`
	Ulid   string `json:"ulid,omitempty"`
	Invite bool   `json:"invite,omitempty"`
}

type EmailAddresses []EmailAddress

func (e *EmailAddresses) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, e)
}

func (e EmailAddresses) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *EmailAddresses) UnmarshalJSON(data []byte) error {
	var full []EmailAddress
	if err := json.Unmarshal(data, &full); err == nil {
		*e = full
		return nil
	}

	var simple []string
	if err := json.Unmarshal(data, &simple); err == nil {
		var result []EmailAddress
		for _, s := range simple {
			result = append(result, EmailAddress{Ulid: s})
		}
		*e = result
		return nil
	}

	return fmt.Errorf("invalid format for EmailAddresses: %s", string(data))
}
