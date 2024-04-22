package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func StartDBConnection(connectionStr string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", connectionStr)
	return
}

func SetFunctionReady(db *sql.DB, id string) (err error) {
	_, err = db.Exec("UPDATE functions SET ready = true WHERE id = $1", id)
	return
}
