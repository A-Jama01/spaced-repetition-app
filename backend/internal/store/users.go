package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID int64 `json:"id"`
	GoogleID string `json:"google_id"`
	Username string	`json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (userStore *UsersStore) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT into users (google_id, username)
		VALUES ($1, $2) RETURNING id, created_at
	`
	err := userStore.db.QueryRowContext(
		ctx,
		query,
		user.GoogleID,	
		user.Username,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
