package models

import "github.com/google/uuid"

// UserBalance struct, that represents user's balance
type UserBalance struct {
	UserId uuid.UUID `json:"userId"`
	Value  float64   `json:"value"`
}
