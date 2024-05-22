package main

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func StartDBConnection(connectionStr string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", connectionStr)
	return
}

func SetFunctionReady(db *sql.DB, id string) (err error) {
	timestamp := time.Now().Unix()
	_, err = db.Exec("UPDATE functions SET built = true, build_timestamp = $2 WHERE id = $1", id, timestamp)
	return
}
