package main

import (
	"fmt"
	"net/http"
)

func (app *app) listDecksHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "list all decks")
}

func (app *app) createDeck(w http.ResponseWriter, r *http.Request) {
}

func (app *app) showDeckHandler(w http.ResponseWriter, r *http.Request) {
}

func (app *app) showDueCardsHandler(w http.ResponseWriter, r *http.Request) {
}

func (app *app) deleteDeckHandler(w http.ResponseWriter, r *http.Request) {
}

func (app *app) updateDeckHandler(w http.ResponseWriter, r *http.Request) {
}
