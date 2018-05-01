package testing

import (
	"github.com/sepetrov/prepaidcard/pkg/internal/model"
	"github.com/sepetrov/prepaidcard/pkg/internal/service/createcard"
)

// Repository is a test helper, which implaments interfaces for interaction
// with the model.
type Repository struct {
	Card *model.Card
	Err  error
}

var _ createcard.Saver = &Repository{}

// SaveCard implements createcard.Saver.
func (r *Repository) SaveCard(card *model.Card) error {
	r.Card = card
	return r.Err
}
