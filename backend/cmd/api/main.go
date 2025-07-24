package main

import (
	"log"

	"github.com/A-Jama01/spaced-repetition-app/internal/db"
	"github.com/A-Jama01/spaced-repetition-app/internal/env"
	"github.com/A-Jama01/spaced-repetition-app/internal/store"
)

func main() {
	cfg := config {
		addr: env.GetString("ADDR", ":3001"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/spaced-repetition?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	log.Println("Database connected")

	store := store.NewStorage(db)

	app := &app {
		config: cfg,
		store: store,
	}

	log.Fatal(app.run(app.routes()))
}
