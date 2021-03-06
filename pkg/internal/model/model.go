// Package model represents the model of the API.
package model

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gofrs/uuid"
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

// NewAuthorizationRequest creates new AuthorizationRequest and blocks amount on card if the request is authorised.
// It returns an error if the request is not authorised.
func NewAuthorizationRequest(card *Card, merchant uuid.UUID, amount uint64) (*AuthorizationRequest, error) {
	if err := card.blockMoney(amount); err != nil {
		return &AuthorizationRequest{}, fmt.Errorf("cannot block the requested amount; %v", err)
	}
	id1, err := uuid.NewV4()
	if err != nil {
		return &AuthorizationRequest{}, fmt.Errorf("cannot generate identifier; %v", err)
	}
	id2, err := uuid.NewV4()
	if err != nil {
		return &AuthorizationRequest{}, fmt.Errorf("cannot generate identifier; %v", err)
	}
	req := &AuthorizationRequest{
		uuid:          id1,
		cardUUID:      card.UUID(),
		merchantUUID:  merchant,
		blockedAmount: amount,
		history: []AuthorizationRequestSnapshot{
			{
				uuid:          id2,
				blockedAmount: amount,
				createdAt:     time.Now(),
			},
		},
	}
	return req, nil
}

// Reverse decreases the blocked amount on card and updates req. It returns error if the request is not authorized.
func (req *AuthorizationRequest) Reverse(card *Card, amount uint64) error {
	if card.UUID() != req.cardUUID {
		return errors.New("cannot reverse from different card")
	}
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	if amount > req.blockedAmount {
		return errors.New("cannot reverse more than the blocked amount")
	}
	if err := card.releaseMoney(amount); err != nil {
		return fmt.Errorf("cannot reverse authorization request; %v", err)
	}
	id, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("cannot generate identifier; %v", err)
	}
	req.blockedAmount -= amount
	req.history = append(
		req.history,
		AuthorizationRequestSnapshot{
			uuid:          id,
			blockedAmount: req.blockedAmount,
			createdAt:     time.Now(),
		},
	)
	return nil
}

// UUID returns the UUID.
func (req *AuthorizationRequest) UUID() uuid.UUID {
	return req.uuid
}

// CardUUID returns the card UUID.
func (req *AuthorizationRequest) CardUUID() uuid.UUID {
	return req.cardUUID
}

// MerchantUUID returns the merchant UUID.
func (req *AuthorizationRequest) MerchantUUID() uuid.UUID {
	return req.merchantUUID
}

// BlockedAmount returns the blocked amount.
func (req *AuthorizationRequest) BlockedAmount() uint64 {
	return req.blockedAmount
}

// CapturedAmount returns the blocked amount.
func (req *AuthorizationRequest) CapturedAmount() uint64 {
	return req.capturedAmount
}

// RefundedAmount returns the blocked amount.
func (req *AuthorizationRequest) RefundedAmount() uint64 {
	return req.refundedAmount
}

// History returns the log of changes.
func (req *AuthorizationRequest) History() []AuthorizationRequestSnapshot {
	return req.history
}

// AuthorizationRequestSnapshot represents a snapshot of AuthorizationRequest.
type AuthorizationRequestSnapshot struct {
	uuid           uuid.UUID
	blockedAmount  uint64
	capturedAmount uint64
	refundedAmount uint64
	createdAt      time.Time
}

// UUID returns the UUID.
func (s AuthorizationRequestSnapshot) UUID() uuid.UUID {
	return s.uuid
}

// BlockedAmount returns the blocked amount.
func (s AuthorizationRequestSnapshot) BlockedAmount() uint64 {
	return s.blockedAmount
}

// CapturedAmount returns the blocked amount.
func (s AuthorizationRequestSnapshot) CapturedAmount() uint64 {
	return s.capturedAmount
}

// RefundedAmount returns the blocked amount.
func (s AuthorizationRequestSnapshot) RefundedAmount() uint64 {
	return s.refundedAmount
}

// CreatedAt returns the time when the snapshot was taken.
func (s AuthorizationRequestSnapshot) CreatedAt() time.Time {
	return s.createdAt
}

// CardData is an interface providing card data.
type CardData interface {
	UUID() uuid.UUID
	AvailableBalance() uint64
	BlockedBalance() uint64
}

// Card represents a prepaid card.
type Card struct {
	uuid             uuid.UUID
	availableBalance uint64
	blockedBalance   uint64
}

// NewCard returns new Card.
func NewCard() (*Card, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("cannot generate identifier; %v", err)
	}
	return &Card{uuid: id}, nil
}

// CardFromData reconstructs card from data.
func CardFromData(data CardData) *Card {
	return &Card{
		uuid:             data.UUID(),
		availableBalance: data.AvailableBalance(),
		blockedBalance:   data.BlockedBalance(),
	}
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

// blockMoney blocks amount from the available balance of c.
func (c *Card) blockMoney(amount uint64) error {
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

// releaseMoney releases blocked amount.
func (c *Card) releaseMoney(amount uint64) error {
	if amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	if amount > c.blockedBalance {
		return errors.New("cannot release more than the blocked balance")
	}
	if amount > math.MaxUint64-c.availableBalance {
		return errors.New("available balance cannot exceed math.MaxUint64")
	}
	c.availableBalance += amount
	c.blockedBalance -= amount
	return nil
}

// chargeMoney reduces the blocked balance with amount.
func (c *Card) chargeMoney(amount uint64) error {
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
