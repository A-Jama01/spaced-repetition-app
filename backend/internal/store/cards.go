package store

import (
	"context"
	"database/sql"
	"time"
)

type Card struct {
	ID int64 `json:"id"`
	DeckID int64 `json:"deck_id"`
	Front string `json:"front"`
	Back string	 `json:"back"`
	Retrievability float64  `json:"retrievability"`
	Stability float64 `json:"stability"`
	Difficulty float64 `json:"difficulty"`
	Due *time.Time `json:"due"`
	CreatedAt time.Time `json:"created_at"`
	LastReview *time.Time `json:"last_review"` 
}

type CardsStore struct {
	db *sql.DB
}

type DueForecast struct {
	DueDate time.Time `json:"due_date"`
	DueCount int64 `json:"due_count"`
}

func (cardStore *CardsStore) Create(ctx context.Context, card *Card) error {
	query := `
	INSERT INTO cards (deck_id, front, back) 
	VALUES ($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()
	
	err := cardStore.db.QueryRowContext(
		ctx,
		query,
		card.DeckID,
		card.Front,
		card.Back,
	).Scan(
		&card.ID,
		&card.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *CardsStore) ListByDeck(ctx context.Context, deckID int64) ([]*Card, error) {
	query := `
	SELECT id, deck_id, front, back, due, last_review, created_at FROM cards
	WHERE deck_id = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*Card
	for rows.Next() {
		var card Card
		err := rows.Scan(
			&card.ID, 
			&card.DeckID, 
			&card.Front, 
			&card.Back, 
			&card.Due, 
			&card.LastReview,
			&card.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		cards = append(cards, &card)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (s *CardsStore) ListDueCards(ctx context.Context, deckID int64) ([]*Card, error) {
	query := `
	SELECT id, deck_id, front, back, due, last_review, created_at FROM cards
	WHERE deck_id = $1 AND (due <= NOW() OR due IS NULL)
	ORDER BY due ASC NULLS LAST, created_at ASC`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*Card
	for rows.Next() {
		var card Card
		err := rows.Scan(
			&card.ID, 
			&card.DeckID, 
			&card.Front, 
			&card.Back, 
			&card.Due, 
			&card.LastReview,
			&card.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		cards = append(cards, &card)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (s *CardsStore) Delete(ctx context.Context, cardID, deckID int64) error {
	query := `
	DELETE FROM cards
	WHERE id = $1 AND deck_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	results, err := s.db.ExecContext(ctx, query, cardID, deckID)
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

func (s *CardsStore) Get(ctx context.Context, cardID, deckID int64) (*Card, error) {
	query := `
	SELECT id, deck_id, front, back, retrievability, stability, difficulty, due, last_review, 
	created_at FROM cards
	WHERE id = $1 AND deck_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var card Card

	err := s.db.QueryRowContext(ctx, query, cardID, deckID).Scan(
		&card.ID,
		&card.DeckID,
		&card.Front,
		&card.Back,
		&card.Retrievability,
		&card.Stability,
		&card.Difficulty,
		&card.Due,
		&card.LastReview,
		&card.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (s *CardsStore) Update(ctx context.Context, card *Card) error {
	query := `
	UPDATE cards
	SET front = $1, 
	back = $2, 
	retrievability = $3, 
	stability = $4, 
	difficulty = $5, 
	due = $6,
	last_review = $7
	WHERE id = $8`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	results, err := s.db.ExecContext(ctx, 
		query, 
		card.Front, 
		card.Back, 
		card.Retrievability,
		card.Stability,
		card.Difficulty,
		card.Due,
		card.LastReview,
		card.ID,
	)
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

func (s *CardsStore) GetDueForecast(ctx context.Context, p StatsQueryParams) ([]*DueForecast, error) {
	query := `
	SELECT DATE(c.due AT TIME ZONE $1) AS due_date, COUNT(*) AS count FROM cards c
	JOIN decks d ON c.deck_id = d.id
	JOIN users u ON d.user_id = u.id
	WHERE u.id = $2
	AND ($3 = '' OR d.name = $3)
	AND DATE(c.due AT TIME ZONE $1) < DATE(NOW() AT TIME ZONE $1) + interval '1 month'
	AND DATE(c.due AT TIME ZONE $1) >= DATE(NOW() AT TIME ZONE $1)
	GROUP BY due_date
	ORDER BY due_date`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, p.TimeZone, p.UserID, p.DeckName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forecasts []*DueForecast
	for rows.Next() {
		var forecast DueForecast

		err = rows.Scan(&forecast.DueDate, &forecast.DueCount)
		if err != nil {
			return nil, err
		}

		forecasts = append(forecasts, &forecast)
	}

	return forecasts, nil
}
