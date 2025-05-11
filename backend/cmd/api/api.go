package main

import (
	"log"
	"net/http"
	"time"
)

type app struct {
	config config
}

type config struct {
	addr string
}

func (app *app) mount() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func (app *app) run(mux *http.ServeMux) error {
	server := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	log.Printf("Server is running on port %s", app.config.addr)
	return server.ListenAndServe()
}
