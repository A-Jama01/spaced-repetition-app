package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUsername = errors.New("Duplicate username")
)

type User struct {
	ID int64 `json:"id"`
	Username string	`json:"username"`
	Password password `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type password struct {
	plaintext *string
	hash []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT into users (username, password_hash)
	VALUES ($1, $2) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password.hash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s UsersStore) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
	SELECT id, username, password_hash, created_at
	FROM users
	WHERE username = $1`	

	var user User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password.hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
