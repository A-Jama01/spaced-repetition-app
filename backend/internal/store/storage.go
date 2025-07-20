package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Decks interface {
		Create(context.Context, *Deck) error
	}
	Cards interface {
		Create(context.Context, *Card) error
	}
	Logs interface {
		Create(context.Context, *Logs) error
	}
	Users interface {
		Create(context.Context, *User) error
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
