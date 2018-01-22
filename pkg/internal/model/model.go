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
	uuid           uuid.UUID
	cardUUID       uuid.UUID
	merchantUUID   uuid.UUID
	blockedAmount  uint64
	capturedAmount uint64
	refundedAmount uint64
	history        []AuthorizationRequestSnapshot
}

// NewAuthorizationRequest creates new AuthorizationRequest if the request is authorised. It returns an error if the
// request is not authorised.
func NewAuthorizationRequest(card Card, merchant uuid.UUID, amount uint64) (*AuthorizationRequest, error) {
	if card.AvailableBalance() < amount {
		return &AuthorizationRequest{}, errors.New("available balance is too low")
	}
	if card.BlockedBalance() > math.MaxUint64-amount {
		return &AuthorizationRequest{}, errors.New("blocked balance cannot exceed math.MaxUint64")
	}
	snapshot := AuthorizationRequestSnapshot{
		uuid:          uuid.NewV4(),
		blockedAmount: amount,
		createdAt:     time.Now(),
	}
	return &AuthorizationRequest{
		uuid:          uuid.NewV4(),
		cardUUID:      card.UUID(),
		merchantUUID:  merchant,
		blockedAmount: amount,
		history:       []AuthorizationRequestSnapshot{snapshot},
	}, nil
}

// AuthorizationRequestSnapshot represents a snapshot of AuthorizationRequest.
type AuthorizationRequestSnapshot struct {
	uuid           uuid.UUID
	blockedAmount  uint64
	capturedAmount uint64
	refundedAmount uint64
	createdAt      time.Time
}

// Card represents a prepaid card.
type Card struct {
	uuid             uuid.UUID
	availableBalance uint64
	blockedBalance   uint64
}

// NewCard returns new Card.
func NewCard() *Card {
	return &Card{uuid: uuid.NewV4()}
}

// UUID returns the UUID.
func (c *Card) UUID() uuid.UUID {
	return c.uuid
}

// AvailableBalance returns the available balance.
func (c *Card) AvailableBalance() uint64 {
	return c.availableBalance
}

// BlockedBalance returns the blocked balance.
func (c *Card) BlockedBalance() uint64 {
	return c.blockedBalance
}

// LoadMoney loads amount onto c.
func (c *Card) LoadMoney(amount uint64) error {
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	if c.availableBalance > math.MaxUint64-amount {
		return errors.New("available balance cannot exceed math.MaxUint64")
	}
	c.availableBalance += amount
	return nil
}

// BlockMoney blocks amount from the available balance of c.
func (c *Card) BlockMoney(amount uint64) error {
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	if amount > c.availableBalance {
		return errors.New("available balance is too low")
	}
	c.availableBalance -= amount
	c.blockedBalance += amount
	return nil
}

// ChargeMoney reduces the blocked balance with amount.
func (c *Card) ChargeMoney(amount uint64) error {
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	if amount > c.blockedBalance {
		return errors.New("blocked balance is too low")
	}
	c.blockedBalance -= amount
	return nil
}

// Transaction represents a transaction associated with a card.
type Transaction struct {
	uuid             uuid.UUID
	cardUUID         uuid.UUID
	eventUUID        uuid.UUID
	eventType        string
	date             time.Time
	amount           uint64
	availableBalance uint64
	blockedBalance   uint64
	description      string
}
