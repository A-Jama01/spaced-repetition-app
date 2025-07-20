package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID int64 `json:"id"`
	Username string	`json:"username"`
	Password string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (userStore *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT into users (username, password)
		VALUES ($1, $2) RETURNING id, created_at
	`
	err := userStore.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
