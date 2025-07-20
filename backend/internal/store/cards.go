package store

import (
	"context"
	"database/sql"
	"time"
	"github.com/lib/pq"
)

type Card struct {
	ID int64 `json:"id"`
	DeckID int64 `json:"deck_id"`
	Front string `json:"front"`
	Back string	 `json:"back"`
	Retrievability float64 `json:"retrievability"`
	Stability float64 `json:"stability"`
	Difficulty float64 `json:"difficulty"`
	Due time.Time `json:"due"`
	Weights []float64 `json:"weights"`
}

type CardsStore struct {
	db *sql.DB
}

func (cardStore *CardsStore) Create(ctx context.Context, card *Card) error {
	query := `
		INSERT INTO cards (deck_id, front, back, retrievability, stability, difficulty, due, 
		weights) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`
	err := cardStore.db.QueryRowContext(
		ctx,
		query,
		card.DeckID,
		card.Front,
		card.Back,
		card.Retrievability,
		card.Stability,
		card.Difficulty,
		card.Due,
		pq.Array(card.Weights),
	).Scan(
		&card.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
