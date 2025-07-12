package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Decks interface {
		CreateDeck(context.Context, *Deck) error
	}
	Cards interface {
		CreateCard(context.Context, *Card) error
	}
	Logs interface {
		CreateLog(context.Context, *Logs) error
	}
	Users interface {
		CreateUser(context.Context, *User) error
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
