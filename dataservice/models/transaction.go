package models

import (
	"github.com/google/uuid"
	"time"
)

// Transaction struct, that represents the transaction
type Transaction struct {
	SenderId    uuid.NullUUID `json:"senderId"`
	ReceiverId  uuid.NullUUID `json:"receiverId"`
	Amount      float64       `json:"amount"`
	Description string        `json:"description"`
	Date        time.Time     `json:"date"`
}
