package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("Record not found")
)

type Storage struct {
	Decks interface {
		Create(context.Context, *Deck) error
		GetByDeckID(context.Context, int64) (*Deck, error)
		GetByUserID(context.Context, int64) ([]*Deck, error)
		DeleteByDeckID(context.Context, int64) error
		Update(context.Context, *Deck) error
	}
	Cards interface {
		Create(context.Context, *Card) error
	}
	Logs interface {
		Create(context.Context, *Logs) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetByUsername(context.Context, string) (*User, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Decks: &DecksStore{db},
		Cards: &CardsStore{db},
		Logs: &LogsStore{db},
		Users: &UsersStore{db},
	}
}
