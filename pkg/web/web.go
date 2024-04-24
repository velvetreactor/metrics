package web

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Web struct {
	db databaser
}

type databaser interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Close(context.Context) error
}

func New() *Web {
	return &Web{}
}

func (w *Web) getDbConn() databaser {
	if w.db != nil {
		return w.db
	}

	dbConn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("unable to connect to database: %v", err)
		os.Exit(1)
	}

	return dbConn
}
