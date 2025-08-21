package main

import (
	"net/http"
	"strconv"
	"github.com/A-Jama01/spaced-repetition-app/internal/scheduler"
	"github.com/A-Jama01/spaced-repetition-app/internal/store"
	"github.com/go-chi/chi/v5"
)

type CardInput struct {
	Front string `json:"front" validate:"required,max=300"`
	Back string	 `json:"back" validate:"required,max=2000"`
}

func (app *app) listCardsHander(w http.ResponseWriter, r *http.Request) {
	deckIDParam := chi.URLParam(r, "deck_id")
	deckID, err := strconv.ParseInt(deckIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	cards, err := app.store.Cards.ListByDeck(ctx, deckID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	
	err = app.writeJSON(w, http.StatusOK, envelope{"cards": cards}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *app) listDueCardsHandler(w http.ResponseWriter, r *http.Request) {
	deckIDParam := chi.URLParam(r, "deck_id")
	deckID, err := strconv.ParseInt(deckIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	cards, err := app.store.Cards.ListDueCards(ctx, deckID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cards": cards}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *app) createCardHandler(w http.ResponseWriter, r *http.Request) {
	var input CardInput
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

	deckIDParam := chi.URLParam(r, "deck_id")
	deckID, err := strconv.ParseInt(deckIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	card := &store.Card{
		DeckID: deckID,
		Front: input.Front,
		Back: input.Back,
	}
	
	ctx := r.Context()
	err = app.store.Cards.Create(ctx, card)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"card": card}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *app) deleteCardHandler(w http.ResponseWriter, r *http.Request) {
	deckIDParam := chi.URLParam(r, "deck_id")
	deckID, err := strconv.ParseInt(deckIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cardIDParam := chi.URLParam(r, "card_id")
	cardID, err := strconv.ParseInt(cardIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	err = app.store.Cards.Delete(ctx, cardID, deckID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Card succesfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *app) updateCardHandler(w http.ResponseWriter, r *http.Request) {
	var input CardInput
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

	deckIDParam := chi.URLParam(r, "deck_id")
	deckID, err := strconv.ParseInt(deckIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cardIDParam := chi.URLParam(r, "card_id")
	cardID, err := strconv.ParseInt(cardIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	
	ctx := r.Context()
	card, err := app.store.Cards.Get(ctx, cardID, deckID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	card.Front = input.Front
	card.Back = input.Back

	err = app.store.Cards.Update(ctx, card)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"card": card}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *app) reviewCardHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Grade int64  `json:"grade" validate:"required,min=1,max=4"`	
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

	deckIDParam := chi.URLParam(r, "deck_id")
	deckID, err := strconv.ParseInt(deckIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cardIDParam := chi.URLParam(r, "card_id")
	cardID, err := strconv.ParseInt(cardIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	card, err := app.store.Cards.Get(ctx, cardID, deckID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	
	err = scheduler.ScheduleCard(card, input.Grade) 
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.store.Cards.Update(ctx, card)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"card": card}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
