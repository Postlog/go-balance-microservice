package utils

import "github.com/google/uuid"

func ToNullableUUID(id uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: id, Valid: true}
}
