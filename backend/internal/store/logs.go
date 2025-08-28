package store

import (
	"context"
	"database/sql"
	"time"
)

type Logs struct {
	ID int64 `json:"id"`
    CardID int64 `json:"card_id"`
    Grade int64 `json:"grade"`
    ReviewedAt time.Time `json:"reviewed_at"`
}

type LogsStore struct {
	db *sql.DB
}

type StatsQueryParams struct {
	UserID int64
	DeckName string
	TimeZone string
}

type ReviewCell struct {
	ReviewDate time.Time `json:"review_date"`
	Reviews int64 `json:"reviews"`
}

func (s *LogsStore) Create(ctx context.Context, logs *Logs) error {
	query := `
		INSERT INTO logs (card_id, grade)
		VALUES ($1, $2) RETURNING id, reviewed_at`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		logs.CardID,
		logs.Grade,
	).Scan(
		&logs.ID,
		&logs.ReviewedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s LogsStore) GetCount(ctx context.Context, p StatsQueryParams) (int64, error) {
	query := `
	SELECT COUNT(*) FROM logs l 
	JOIN cards c ON l.card_id = c.id 
	JOIN decks d ON c.deck_id = d.id
	JOIN users u ON d.user_id = u.id
	WHERE u.id = $1 
	AND ($2 = '' OR d.name = $2) 
	AND DATE(l.reviewed_at AT TIME ZONE $3) = DATE(NOW() AT TIME ZONE $3)`
	
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	var count int64	
	err := s.db.QueryRowContext(ctx, query, p.UserID, p.DeckName, p.TimeZone).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *LogsStore) GetRetention(ctx context.Context, p StatsQueryParams) (float64, error) {
	query := `
	SELECT COALESCE(
		ROUND(
			COUNT(*) FILTER (WHERE l.grade > 1)::numeric / NULLIF(COUNT(*), 0), 
			2) * 100, 
		0) AS accuracy
	FROM logs l
	JOIN cards c ON l.card_id = c.id
	JOIN decks d ON c.deck_id = d.id
	JOIN users u ON d.user_id = u.id
	WHERE u.id = $1
	AND ($2 = '' OR d.name = $2)
	AND DATE(l.reviewed_at AT TIME ZONE $3) = DATE(NOW() AT TIME ZONE $3)`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var accuracy float64
	err := s.db.QueryRowContext(ctx, query, p.UserID, p.DeckName, p.TimeZone).Scan(&accuracy)
	if err != nil {
		return 0.0, err
	}

	return accuracy, nil
}

func (s *LogsStore) GetHeatMap(ctx context.Context, p StatsQueryParams) ([]*ReviewCell, error) {
	query := `
	SELECT DATE(l.reviewed_at AT TIME ZONE $1) AS date, COUNT(*) AS reviews FROM logs l
	JOIN cards c ON l.card_id = c.id
	JOIN decks d ON c.deck_id = d.id
	JOIN users u ON d.user_id = u.id
	WHERE u.id = $2
	AND ($3 = '' OR d.name = $3)
	AND l.reviewed_at >= NOW() - INTERVAL '1 year'
	GROUP BY date
	ORDER BY date` 

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, p.TimeZone, p.UserID, p.DeckName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var heatmap []*ReviewCell
	for rows.Next() {
		var cell ReviewCell

		err := rows.Scan(&cell.ReviewDate, &cell.Reviews)
		if err != nil {
			return nil, err
		}
		
		heatmap = append(heatmap, &cell)
	}

	return heatmap, nil
}
