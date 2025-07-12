package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID int64
	GoogleID string			
	Username string	
	CreatedAt time.Time
}

type UsersStore struct {
	db *sql.DB
}

func (userStore *UsersStore) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT into users (google_id, username, created_at)
		VALUES ($1, $2, $3) RETURNING id
	`
	err := userStore.db.QueryRowContext(
		ctx,
		query,
		user.GoogleID,	
		user.Username,
		user.CreatedAt,
	).Scan(
		&user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
