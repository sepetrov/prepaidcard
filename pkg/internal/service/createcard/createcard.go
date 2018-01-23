package createcard

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/sepetrov/prepaidcard/pkg/internal/event"
	"github.com/sepetrov/prepaidcard/pkg/internal/model"
)

// Response is the response, which Service returns when a card is successfully created.
type Response struct {
	UUID             string `json:"uuid"`
	AvailableBalance string `json:"availableBalance"`
	BlockedBalance   string `json:"blockedBalance"`
}

// Service is the service creating new cards.
type Service struct {
	saver      Saver
	dispatcher Dispatcher
}

// New returns new service creating cards.
func New(s Saver, d Dispatcher) *Service {
	return &Service{s, d}
}

// CreateCard creates a new card.
func (svc *Service) CreateCard() (Response, error) {
	card := model.NewCard()
	if err := svc.saver.SaveCard(card); err != nil {
		return Response{}, fmt.Errorf("CreateCard() cannot persist card; %v", err)
	}
	svc.dispatcher.DispatchCardCreated(event.CardCreated{
		UUID:     uuid.NewV4(),
		Time:     time.Now(),
		CardUUID: card.UUID(),
	})
	return Response{
		UUID:             card.UUID().String(),
		AvailableBalance: string(card.AvailableBalance()),
		BlockedBalance:   string(card.BlockedBalance()),
	}, nil
}

// Saver is interface for persistence of new cards.
type Saver interface {
	SaveCard(*model.Card) error
}

// Dispatcher is an interface for dispatching CardCreated event.
type Dispatcher interface {
	DispatchCardCreated(event.CardCreated)
}
