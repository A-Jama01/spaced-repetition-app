package main

import (
	"log"
	"net/http"
	"time"
	"github.com/A-Jama01/spaced-repetition-app/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	//"github.com/go-chi/jwtauth"
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


func (app *app) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			//Public routes
			r.Group(func(r chi.Router) {
				r.Route("/auth", func(r chi.Router){
					r.Post("/register", app.registerHandler)	
					r.Post("/login", app.loginHandler)	
					r.Post("/logout", app.logoutHandler)
				})
			})

			//Protected routes
			r.Group(func(r chi.Router) {
				r.Route("/decks", func(r chi.Router) {
					r.Get("/", app.listDecksHandler)
					r.Post("/", app.createDeck)
					r.Route("/{deck_id}", func(r chi.Router) {
						r.Get("/", app.showDeckHandler)
						r.Get("/due", app.showDueCardsHandler)
						r.Delete("/", app.deleteDeckHandler)
						r.Put("/", app.updateDeckHandler)
					})
				})

				r.Route("/cards", func(r chi.Router) {
					r.Get("/", app.listCardsHander)
					r.Post("/", app.createCardHandler)
					r.Route("/{card_id}", func(r chi.Router) {
						r.Delete("/", app.deleteCardHandler)
						r.Put("/", app.editCardHandler)
						r.Patch("/review", app.reviewCardHandler)
					})	
				})

				r.Route("/stats", func(r chi.Router) {
					r.Get("/", app.listStatsHandler)
				})
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
