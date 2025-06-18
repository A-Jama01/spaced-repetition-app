package store

import (
	"context"
	"database/sql"
)

type DecksStore struct {
	db *sql.DB
}

func (s *DecksStore) Create(ctx context.Context) error {
	return nil
}
