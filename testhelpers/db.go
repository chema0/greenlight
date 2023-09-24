package testhelpers

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/chema0/greenlight/config"
	_ "github.com/lib/pq"
)

func NewTestingDB() *sql.DB {
	cfg := config.NewConfig("test")

	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	db.SetConnMaxIdleTime(cfg.DB.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: run migrations

	return db
}
