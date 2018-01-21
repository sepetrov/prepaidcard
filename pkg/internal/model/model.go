// Package model represents the model of the API.
package model

import (
	"time"

	"github.com/satori/go.uuid"
)

// AuthorizationRequest represents the requests sent by a merchant to charge a customer.
type AuthorizationRequest struct {
	UUID           uuid.UUID
	CardUUID       uuid.UUID
	MerchantUUID   uuid.UUID
	BlockedAmount  uint64
	CapturedAmount uint64
	RefundedAmount uint64
	History        []AuthorizationRequestSnapshot
}

// AuthorizationRequestSnapshot represents a snapshot of AuthorizationRequest.
type AuthorizationRequestSnapshot struct {
	UUID           uuid.UUID
	BlockedAmount  uint64
	CapturedAmount uint64
	RefundedAmount uint64
	CreatedAt      time.Time
}

// Card represents a prepaid card.
type Card struct {
	UUID             uuid.UUID
	AvailableBalance uint64
	BlockedBalance   uint64
}

// Transaction represents a transaction associated with a card.
type Transaction struct {
	UUID             uuid.UUID
	CardUUID         uuid.UUID
	EventUUID        uuid.UUID
	EventType        string
	Date             time.Time
	Amount           uint64
	AvailableBalance uint64
	BlockedBalance   uint64
	Description      string
}
