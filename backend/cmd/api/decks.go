package main

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/A-Jama01/spaced-repetition-app/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *app) listDecksHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "list all decks")
}

func (app *app) createDeck(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int64 `json:"id" validate:"required"` 
		UserID int64 `json:"user_id" validate:"required"`
		Name string `json:"name" validate:"required"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	
	err = app.validate.Struct(input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *app) showDeckHandler(w http.ResponseWriter, r *http.Request) {
	deckIDParam := chi.URLParam(r, "deck_id")
	deckID, err := strconv.Atoi(deckIDParam)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	
	deck := store.Deck{
		ID:  int64(deckID),
		UserID: 1,
		Name: "Algorithms",
	}

	err = app.writeJSON(w, http.StatusOK, envolope{"deck": deck}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *app) showDueCardsHandler(w http.ResponseWriter, r *http.Request) {
}

func (app *app) deleteDeckHandler(w http.ResponseWriter, r *http.Request) {
}

func (app *app) updateDeckHandler(w http.ResponseWriter, r *http.Request) {
}
