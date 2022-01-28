// Package protocol provides protocol for interaction through HTTP API
package protocol

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/dataservice/models"
)

// Response represents the response to the client
type Response struct {
	Error   *Error          `json:"error"`
	Payload json.RawMessage `json:"payload"`
}

type Error struct {
	Message string `json:"message"`
	Code    *int   `json:"code"`
}

type BalancePayload struct {
	Balance models.UserBalance `json:"balance"`
}

type TransactionsPayload struct {
	Transactions []models.Transaction `json:"transactions"`
}

type ValidatableRequest interface {
	Validate() error
}

// UpdateBalanceRequest represents request from the client, that is going to add to the user's balance or reduce it
//
// Implements ValidatableRequest interface
type UpdateBalanceRequest struct {
	UserId      *uuid.UUID `json:"userId"` // required field
	Amount      *float64   `json:"amount"` // required field
	Currency    string     `json:"currency"`
	Description string     `json:"description"`
}

func (r *UpdateBalanceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.UserId, validation.Required),
		validation.Field(&r.Amount, validation.Required),
	)
}

// TransferFoundsRequest represents request from the client, that is going to transfer founds from one user to another
//
// Implements ValidatableRequest interface
type TransferFoundsRequest struct {
	SenderId    *uuid.UUID `json:"senderId"`   // required field
	ReceiverId  *uuid.UUID `json:"receiverId"` // required field
	Amount      *float64   `json:"amount"`     // required field
	Currency    string     `json:"currency"`
	Description string     `json:"description"`
}

func (r *TransferFoundsRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.SenderId, validation.Required),
		validation.Field(&r.ReceiverId, validation.Required),
		validation.Field(&r.Amount, validation.Required),
	)
}

// GetBalanceRequest represents request from the client, that is going to get user's balance
//
// Implements ValidatableRequest interface
type GetBalanceRequest struct {
	UserId *uuid.UUID `json:"userId"` // required field
}

func (r *GetBalanceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.UserId, validation.Required),
	)
}

// GetTransactionsRequest represents request from the client, that is going to get user's transactions history
//
// Implements ValidatableRequest interface
type GetTransactionsRequest struct {
	UserId         *uuid.UUID `json:"userId"` // required field
	OrderBy        string     `json:"orderBy"`
	OrderDirection string     `json:"orderDirection"`
	Limit          int        `json:"limit"`
	Offset         int        `json:"offset"`
}

func (r *GetTransactionsRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.UserId, validation.Required),
	)
}
