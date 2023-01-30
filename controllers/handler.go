package controllers

import (
	"log"
	"fmt"
	"database/sql"
	"github.com/enzdor/gomsg/sqlc"
)

type Handler struct {
	q *sqlc.Queries
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	queries := sqlc.New(db)

	return &Handler {
		q: queries,
		db: db,
	}
}

func NewDB(user string, pass string, name string) *sql.DB{
	cfg := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", user, pass, name)

	db, err := sql.Open("mysql", cfg)
	if err != nil {
	    log.Fatal(err)
	}

	return db
}













