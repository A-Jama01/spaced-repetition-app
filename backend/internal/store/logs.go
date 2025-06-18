package store

import (
	"context"
	"database/sql"
)

type LogsStore struct {
	db *sql.DB
}

func (s *LogsStore) Create(ctx context.Context) error {
	return nil
}
