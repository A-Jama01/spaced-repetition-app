package store

import (
	"context"
	"database/sql"
	"time"
	"github.com/A-Jama01/spaced-repetition-app/internal"
)

type Logs struct {
    ID int64    
    CardID int64
    ReviewDate time.Time
    Grade internal.Grade
}

type LogsStore struct {
	db *sql.DB
}

func (logsStore *LogsStore) CreateLog(ctx context.Context, logs * Logs) error {
	query := `
		INSERT INTO logs (card_id, review_date, grade)
		VALUES ($1, $2, $3) RETURNING id
	`
	err := logsStore.db.QueryRowContext(
		ctx,
		query,
		logs.CardID,
		logs.ReviewDate,
		logs.Grade,
	).Scan(
		&logs.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
