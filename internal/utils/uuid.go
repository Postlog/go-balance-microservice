package utils

import "github.com/google/uuid"

func ParseUUID(s string) uuid.NullUUID {
	parsed, err := uuid.Parse(s)
	id := uuid.NullUUID{UUID: parsed, Valid: true}
	if err != nil {
		id.Valid = false
	}
	return id
}
