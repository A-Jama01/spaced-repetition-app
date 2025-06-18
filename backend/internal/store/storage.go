package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Decks interface {
		Create(context.Context) error
	}
	Cards interface {
		Create(context.Context) error
	}
	Logs interface {
		Create(context.Context) error
	}
	Users interface {
		Create(context.Context) error
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
