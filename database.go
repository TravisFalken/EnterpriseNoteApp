package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

//Delete Specific note
func deleteSpecificNote(r *http.Request) (noteDeleted bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//get the id of the note the user wants to delete
	noteid := mux.Vars(r)["id"]

	//get the actually username out of the cookie
	username := getUserName(r)
	stmt, err := db.Prepare("DELETE FROM _note WHERE note_owner=$1 AND note_id=$2;")
	if err != nil {
		log.Fatal(err)
	}

	deleted, _ := stmt.Exec(username, noteid)
	rowsaffected, _ := deleted.RowsAffected()
	if rowsaffected > 0 {
		return true
	}
	return false
}
