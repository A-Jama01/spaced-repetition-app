package main

import (
	"log"

	"github.com/A-Jama01/spaced-repetition-app/internal/env"
	"github.com/A-Jama01/spaced-repetition-app/internal/store"
)

func main() {
	cfg := config {
		addr: env.GetString("ADDR", ":3001"),
	}

	store := store.NewStorage(nil)

	app := &app {
		config: cfg,
		store: store,
	}

	log.Fatal(app.run(app.mount()))
}
