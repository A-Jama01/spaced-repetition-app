package main

import (
	"net/http"
	"strconv"

	"github.com/A-Jama01/spaced-repetition-app/internal/store"
	"github.com/go-chi/chi/v5"
)

type DeckInput struct {
	Name string `json:"name" validate:"required,min=1,max=70"`
}

func (app *app) listDecksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `validate:"max=70"`
	}

	queryString := r.URL.Query()
	input.Name = app.readString(queryString, "name", "")

	err := app.validate.Struct(input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	
	userID, err := app.getUserIDFromContext(ctx)	
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	decks, err := app.store.Decks.ListAll(ctx, userID, input.Name)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"decks": decks}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *app) createDeck(w http.ResponseWriter, r *http.Request) {
	var input DeckInput

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

	ctx := r.Context()
	
	userID, err := app.getUserIDFromContext(ctx)	
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	deck := &store.Deck{
		UserID: userID,
		Name: input.Name,
	}
	
	err = app.store.Decks.Create(ctx, deck)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"deck": deck}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *app) deleteDeckHandler(w http.ResponseWriter, r *http.Request) {
	deckIDParam := chi.URLParam(r, "deck_id")
	deckID, err := strconv.ParseInt(deckIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	
	ctx := r.Context()
	err = app.store.Decks.DeleteByDeckID(ctx, int64(deckID))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Deck deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *app) updateDeckHandler(w http.ResponseWriter, r *http.Request) {
	deckIDParam := chi.URLParam(r, "deck_id")
	deckID, err := strconv.ParseInt(deckIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	deck, err := app.store.Decks.GetByDeckID(ctx, deckID)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input DeckInput 

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.validate.Struct(input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	deck.Name = input.Name
	
	err = app.store.Decks.Update(ctx, deck)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"deck": deck}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
