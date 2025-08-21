package store

import (
	"context"
	"database/sql"
	"time"
)

type Deck struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type DecksStore struct {
	db *sql.DB
}

func (s *DecksStore) Create(ctx context.Context, deck *Deck) error {
	query := `
	INSERT INTO decks (user_id, name)
	VALUES ($1, $2) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		deck.UserID,
		deck.Name,
	).Scan(
		&deck.ID,
		&deck.CreatedAt,
	)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *DecksStore) GetByDeckID(ctx context.Context, deckID int64) (*Deck, error) {
	query := `
	SELECT * FROM decks
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	var deck Deck
	
	err := s.db.QueryRowContext(ctx, query, deckID).Scan(
		&deck.ID,
		&deck.UserID,
		&deck.Name,
		&deck.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &deck, nil
}

func (s *DecksStore) ListByUserID(ctx context.Context, userID int64) ([]*Deck, error) {
	query := `
	SELECT * FROM decks 
	WHERE user_id = $1`

	ctx, cancel := 	context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decks []*Deck
	for rows.Next() {
		var deck Deck
		err := rows.Scan(&deck.ID, &deck.UserID, &deck.Name, &deck.CreatedAt)
		if err != nil {
			return nil, err
		}

		decks = append(decks, &deck)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return decks, nil
}

func (s *DecksStore) DeleteByDeckID(ctx context.Context, deckID int64) error {
	if deckID < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM decks
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, deckID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (s *DecksStore) Update(ctx context.Context, deck *Deck) error {
	query := `
	UPDATE decks
	SET name = $1
	WHERE id = $2`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	results, err := s.db.ExecContext(ctx, query, deck.Name, deck.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
