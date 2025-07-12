package store

import (
	"context"
	"database/sql"
)

type Deck struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Name string `json:"name"`
}

type DecksStore struct {
	db *sql.DB
}

func (s *DecksStore) CreateDeck(ctx context.Context, deck *Deck) error {
	query := `
		INSERT INTO decks (user_id, name)
		VALUES ($1, $2) RETURNING id 
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		deck.UserID,
		deck.Name,
	).Scan(
		&deck.ID,
	)
	if err != nil {
		return err
	}
	
	return nil
}
