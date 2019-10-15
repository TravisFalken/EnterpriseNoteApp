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

//===================Validate Username=============================
//Validate if the username already exists in the database  (username has to be unique)
//Return true if username exists
func userNameExists(username string) bool {
	var name string
	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//Prepare query to check if the username already exists
	getUserName, err := db.Prepare("Select user_name FROM _user WHERE user_name = $1;")
	if err != nil {
		log.Fatal(err)
	}
	err = getUserName.QueryRow(username).Scan(&name)
	//if error username does not exist in database
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	//Username does exist in the database
	return true
}

//===================Validate Password=============================
func validatePass(password string, username string) bool {
	var pass string
	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//prepare statement to check for password
	stmt, err := db.Prepare("SELECT password FROM _user WHERE password = $1 AND user_name = $2;")
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(password, username).Scan(&pass)
	//if nothing is returned
	if err == sql.ErrNoRows {
		//password does not match
		return false
	}
	if err != nil {
		log.Fatal(err)
	}

	//password matches
	return true

}

//===================Delete Specific note=============================

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

//===================Add User=============================
//
func addUserSQL(newUser User) string {

	db := connectDatabase()
	defer db.Close()

	if !userNameExists(newUser.UserName) {
		//Prepare insert to stop SQL injections
		log.Println("Entered add user if statement")
		stmt, err := db.Prepare("INSERT INTO _user VALUES($1,$2,$3,$4,$5);")
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(newUser.UserName, newUser.Password, newUser.Email, newUser.GivenName, newUser.FamilyName)
		if err != nil {
			log.Fatal(err)
		}
		return "Added user"
	}

	return "Username already exists"

}

//===================Add Note=============================
//

func addNoteSQL(newNote Note) string {

	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//Prepare insert to stop SQL injections
	stmt, err := db.Prepare("INSERT INTO _note (title, body, date_created, note_owner) VALUES($1,$2,$3,$4);")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(newNote.NoteTitle, newNote.NoteBody, newNote.CreatedDate, newNote.NoteOwner)
	if err != nil {
		log.Fatal(err)
	}

	return "New Note Added"
	//User is not logged in

}

//===================List All Notes=============================
//

func listAllNotesSQL(username string) []Note {

	var notes []Note
	//Connect to db
	db := connectDatabase()
	defer db.Close()
	stmt, err := db.Prepare("SELECT _note.note_id, _note.note_owner, _note.title, _note.body, _note.date_created FROM _note LEFT OUTER JOIN _note_privileges ON (_note.note_id = _note_privileges.note_id) WHERE _note.note_owner = $1 OR _note_privileges.user_name = $1;")
	var note Note
	rows, err := stmt.Query(username)
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.NoteOwner, &note.NoteTitle, &note.NoteBody, &note.CreatedDate)
		if err != nil {
			log.Fatal(err)
		}
		notes = append(notes, note)

	}
	return notes
}

//===================List All Users=============================
//
func listAllUsersSQL(loggedInUser string) []string {
	var users []string
	var username string
	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//Send query to the db
	rows, err := db.Query("SELECT user_name FROM _user;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&username)
		//Make sure we dont add logged in user to recommended
		if loggedInUser != username {
			users = append(users, username)
		}
	}
	return users
}

//===================Partial Text Search=============================

func partialTextSearchSQL(bodyText string, username string) []Note {
	var notes []Note
	//Connect to db
	db := connectDatabase()
	defer db.Close()

	bodyText += ":*" //for testing
	stmt, err := db.Prepare("SELECT _note.note_id, _note.note_owner, _note.title, _note.body, _note.date_created FROM _note LEFT OUTER JOIN _note_privileges ON (_note.note_id = _note_privileges.note_id) WHERE body ~ $2 AND _note.note_owner = $1 OR _note_privileges.user_name = $1;")
	if err != nil {
		log.Fatal(err)
	}
	var note Note
	rows, err := stmt.Query(username, bodyText)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody, &note.CreatedDate, &note.NoteOwner)
		if err != nil {
			log.Fatal(err)
		}
		notes = append(notes, note)

	}
	return notes
}
