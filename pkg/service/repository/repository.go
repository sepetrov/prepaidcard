// Package repository has the repository service, which provides interface with
// persistance layer of the API.
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid"

	"github.com/sepetrov/prepaidcard/pkg/internal/model"
	"github.com/sepetrov/prepaidcard/pkg/internal/service/createcard"
)

const cardCollection = "card"

const sqlInsertCard = "INSERT INTO card (uuid, available_balance, blocked_balance) VALUES (?, ?, ?)"
const sqlSelectCard = "SELECT uuid, available_balance, blocked_balance FROM card WHERE uuid = ? LIMIT 1"

// ErrNotFound is returned when the excepcted record(s) can not be found.
var ErrNotFound = errors.New("record not found")

// Repository is a service, which provides interface with persistance layer.
type Repository struct {
	db *sql.DB
}

// New returns new repository for db.
func New(db *sql.DB) *Repository {
	return &Repository{db}
}

var _ createcard.Saver = &Repository{}

// card represents card data
type card struct {
	uuid             uuid.UUID
	availableBalance uint64
	blockedBalance   uint64
}

// Ensure card implements model.CardData.
var _ model.CardData = &card{}

// UUID returns the UUID.
func (c card) UUID() uuid.UUID {
	return c.uuid
}

// AvailableBalance returns the available balance.
func (c card) AvailableBalance() uint64 {
	return c.availableBalance
}

// BlockedBalance returns the blocked balance.
func (c card) BlockedBalance() uint64 {
	return c.blockedBalance
}

// SaveCard persits new card.
func (r *Repository) SaveCard(card *model.Card) error {
	stmt, err := r.db.Prepare(sqlInsertCard)
	if err != nil {
		return fmt.Errorf("cannot prepare statement to save card: %v", err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(card.UUID(), card.AvailableBalance(), card.BlockedBalance()); err != nil {
		log.Fatal(err)
	}
	return nil
}

// GetCard returns the card with uuid.
func (r *Repository) GetCard(uuid uuid.UUID) (*model.Card, error) {
	data := card{}
	row := r.db.QueryRow(sqlSelectCard, uuid.String())
	err := row.Scan(&data.uuid, &data.availableBalance, &data.blockedBalance)
	if err == sql.ErrNoRows {
		return &model.Card{}, ErrNotFound
	}
	if err != nil {
		return &model.Card{}, fmt.Errorf("got error, want one row: %v", err)
	}
	return model.CardFromData(data), nil
}
