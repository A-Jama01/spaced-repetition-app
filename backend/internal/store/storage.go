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
		ListAll(context.Context, int64, string) ([]*Deck, error)
		DeleteByDeckID(context.Context, int64) error
		Update(context.Context, *Deck) error
	}
	Cards interface {
		Create(context.Context, *Card) error
		ListByDeck(context.Context, int64, string) ([]*Card, error)
		ListDueCards(context.Context, int64) ([]*Card, error)
		Delete(ctx context.Context, cardID, deckID int64) error
		Get(ctx context.Context, cardID, deckID int64) (*Card, error)
		Update(context.Context, *Card) error
		GetDueForecast(context.Context, StatsQueryParams) ([]*DueForecast, error)
	}
	Logs interface {
		Create(context.Context, *Logs) error
		GetCount(context.Context, StatsQueryParams) (int64, error)
		GetRetention(context.Context, StatsQueryParams) (float64, error)
		GetHeatMap(context.Context, StatsQueryParams) ([]*ReviewCell, error)
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
