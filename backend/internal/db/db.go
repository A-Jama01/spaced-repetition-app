package db

import (
	"context"
	"database/sql"
	"time"
)

func New(dsn string, maxOpenConn int, maxIdleConn int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)

	idleDuration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(idleDuration)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	db.PingContext(ctx)
	
	return db, nil
}
