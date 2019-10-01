package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func connectDatabase() (db *sql.DB) {
	//Open db connection
	db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")

	if err != nil {
		log.Panic(err)
	}

	return db
}
