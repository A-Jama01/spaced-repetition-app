package main

import (
	"log"
	"net/http"
	"time"

	"github.com/A-Jama01/spaced-repetition-app/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type app struct {
	config config
	store store.Storage
}

type config struct {
	addr string
	db dbConfig
}

type dbConfig struct {
	addr string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime string
}

func (app *app) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router){
				r.Post("/register", registerHandler)	
				r.Post("/login", loginHandler)	
			})

			r.Route("/decks", func(r chi.Router) {
				r.Get("/", listDecksHandler)
				r.Post("/", createDeck)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", getDeckHandler)
					r.Get("/due", getDueCardsHandler)
					r.Delete("/", deleteDeckByIDHandler)
					r.Patch("/", updateDeckByIDHandler)
				})
			})
			
			r.Route("/cards", func(r chi.Router) {
				r.Get("/", listCardsHander)
				r.Post("/", createCardHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Delete("/", deleteCardByIDHandler)
					r.Patch("/", updateCardByIDHandler)
					r.Patch("/review", reviewCardByIDHandler)
				})	
			})

			r.Route("/stats", func(r chi.Router) {
				r.Get("/", listStatsHandler)
			})
		})
	})

	return r
}

func (app *app) run(mux http.Handler) error {
	server := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout: time.Minute,
	}

	log.Printf("Server is running on port %s", app.config.addr)
	return server.ListenAndServe()
}
