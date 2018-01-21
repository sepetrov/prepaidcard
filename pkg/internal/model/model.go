// Package model represents the model of the API.
package model

import (
	"errors"
	"math"
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

// LoadMoney loads amount onto c.
func (c *Card) LoadMoney(amount uint64) error {
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	if c.AvailableBalance > math.MaxUint64-amount {
		return errors.New("available balance cannot exceed math.MaxUint64")
	}
	c.AvailableBalance += amount
	return nil
}

// BlockMoney blocks amount from the available balance of c.
func (c *Card) BlockMoney(amount uint64) error {
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	if amount > c.AvailableBalance {
		return errors.New("available balance is too low")
	}
	c.AvailableBalance -= amount
	c.BlockedBalance += amount
	return nil
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
