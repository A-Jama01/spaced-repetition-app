package main

import (
	"log"
	"os"

	"github.com/A-Jama01/spaced-repetition-app/internal/db"
	"github.com/A-Jama01/spaced-repetition-app/internal/env"
	"github.com/A-Jama01/spaced-repetition-app/internal/store"
	"github.com/go-playground/validator/v10"
	"github.com/go-chi/jwtauth/v5"
)

func main() {
	cfg := config {
		addr: env.GetString("ADDR", ":3001"),
		db: dbConfig{
			dsn: env.GetString("SR_DB_DSN", "postgres://admin:adminpassword@localhost/spaced-repetition?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.dsn,
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
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	validate := validator.New(validator.WithRequiredStructEnabled())
	tokenAuth := jwtauth.New("HS256", []byte(env.GetString("JWT_SECRET", "secret")), nil)

	app := &app {
		config: cfg,
		store: store,
		logger: logger,
		validate: validate,
		jwtAuth: tokenAuth,
	}

	err = app.run(app.routes())
	logger.Fatal(err)
}
