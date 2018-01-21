// Package event has the event generated as a result of the requests which the API is handling.
package event

import (
	"time"

	"github.com/satori/go.uuid"
)

// CardCreated represents the registration of a new card to the system.
type CardCreated struct {
	UUID     uuid.UUID
	Time     time.Time
	CardUUID uuid.UUID
}

// CardLoaded represents the loading of a card by the user.
type CardLoaded struct {
	UUID     uuid.UUID
	Time     time.Time
	CardUUID uuid.UUID
	Amount   uint64
}

// AuthorizationRequestCreated represents the submission of an authorization request from a merchant.
type AuthorizationRequestCreated authorizationRequest

// AuthorizationRequestReversed represents the reversal of an authorization request from a merchant.
type AuthorizationRequestReversed authorizationRequest

// AuthorizationRequestCaptured represents the capturing of a transaction by merchant.
type AuthorizationRequestCaptured authorizationRequest

type authorizationRequest struct {
	UUID         uuid.UUID
	Time         time.Time
	CardUUID     uuid.UUID
	MerchantUUID uuid.UUID
}
