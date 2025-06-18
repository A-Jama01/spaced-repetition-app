package store

import (
	"context"
	"database/sql"
)

type CardsStore struct {
	db *sql.DB
}

func (s *CardsStore) Create(ctx context.Context) error {
	return nil
}
