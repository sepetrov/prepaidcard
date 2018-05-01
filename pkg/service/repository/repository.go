// Package repository has the repository service, which provides interface with
// persistance layer of the API.
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/sepetrov/prepaidcard/pkg/internal/model"
	"github.com/sepetrov/prepaidcard/pkg/internal/service/createcard"
)

const cardCollection = "card"

const sqlInsertCard = "INSERT INTO card (uuid, available_balance, blocked_balance) VALUES (?, ?, ?)"

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
