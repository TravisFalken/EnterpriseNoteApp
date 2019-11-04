package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func connectDatabase() (db *sql.DB) {
	// enterprisenoteapp_db_1
	//connString := "postgresql://postgres:password@enterprisenoteapp_db_1:5432?sslmode=disable"
	// connString := "host=db port=5432 user=postgres password=password dbname=noteBookApp sslmode=disable"

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

		_, err = stmt.Exec(newUser.GivenName, newUser.FamilyName, newUser.UserName, newUser.Password, newUser.Email)
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

//===================Partial Text Search in Body=============================

func partialTextBodySearchSQL(bodyText string, username string) []Note {
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

//==============Search owned notes using partial text based on the notes title=================
func partialSeachOwnedTitle(searchText string, r *http.Request) (ownedNotes []Note) {
	//get username
	username := getUserName(r)
	//Connect to database

	db := connectDatabase()
	defer db.Close()
	//searchText += ":*"
	stmt, err := db.Prepare("SELECT _note.note_id, _note.note_owner, _note.title, _note.body, _note.date_created FROM _note WHERE _note.title ~* $2 AND _note.note_owner = $1;")
	if err != nil {
		log.Fatal(err)
	}
	var note Note
	rows, err := stmt.Query(username, searchText)
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.NoteOwner, &note.NoteTitle, &note.NoteBody, &note.CreatedDate)
		if err != nil {
			log.Panic(err)
		}
		ownedNotes = append(ownedNotes, note)
	}
	return ownedNotes
}

//=============Partial search notes you are appart of by their title=================
func partialSearchPartOfTitle(titleText string, r *http.Request) (partOfNotes []Note) {
	//get username
	username := getUserName(r)
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//titleText += ":*"
	stmt, err := db.Prepare("SELECT _note.note_id, _note.note_owner, _note.title, _note.body, _note.date_created FROM _note_privileges JOIN _note ON _note_privileges.note_id = _note.note_id WHERE _note.title ~* $2 AND _note_privileges.user_name = $1")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Got HERE: Part of title") // for testing
	var note Note
	rows, err := stmt.Query(username, titleText)
	if err != nil {
		log.Panic(err)
		return partOfNotes
	}
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.NoteOwner, &note.NoteTitle, &note.NoteBody, &note.CreatedDate)
		if err != nil {
			log.Panic(err)
		}
		partOfNotes = append(partOfNotes, note)
	}

	return partOfNotes
}

//===================Delete Specific note=============================
// This function currently tested as function below

func deleteSpecificNoteSQL(noteid string, username string) (noteDeleted bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//get the actually username out of the cookie

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

//===================Delete All Notes By User=============================

func deleteAllUserNotesSQL(username string) (noteDeleted bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//get the actually username out of the cookie

	stmt, err := db.Prepare("DELETE FROM _note WHERE note_owner=$1;")
	if err != nil {
		log.Fatal(err)
	}

	deleted, _ := stmt.Exec(username)
	rowsaffected, _ := deleted.RowsAffected()
	if rowsaffected > 0 {
		return true
	}
	return false
}

//===================Delete Specific user=============================

func deleteSpecificUserSQL(username string) (noteDeleted bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//get the actually username out of the cookie

	stmt, err := db.Prepare("DELETE FROM _user WHERE user_name=$1;")
	if err != nil {
		log.Fatal(err)
	}

	deleted, _ := stmt.Exec(username)
	rowsaffected, _ := deleted.RowsAffected()
	if rowsaffected > 0 {
		return true
	}
	return false
}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR HTTP TO CONTROLLER
// ------------------------------------------------------------------------------------------------------------------------------------------------

//check if user has read permissions
func readPermissions(r *http.Request) (readPremission bool) {
	username := getUserName(r)
	var read string
	//Get note id from http
	noteid := mux.Vars(r)["id"]
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	stmt, err := db.Prepare("SELECT read FROM _note_privileges WHERE user_name = $1 AND note_id = $2")
	//if no rows were returned
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Panic(err)
	}
	err = stmt.QueryRow(username, noteid).Scan(&read)
	if err != nil {
		readPremission = false
		return readPremission
	}
	if read == "t" {
		readPremission = true
		return readPremission
	}

	readPremission = false
	return readPremission

}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR HTTP TO CONTROLLER
// ------------------------------------------------------------------------------------------------------------------------------------------------

//Check if user has write permissions
func checkWritePermissions(r *http.Request) (writePermission bool) {
	username := getUserName(r)
	var write string

	noteid := mux.Vars(r)["id"]
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//prepare statement
	stmt, err := db.Prepare("SELECT write FROM _note_privileges WHERE user_name = $1 AND note_id = $2")
	if err != nil {
		log.Panic(err)
	}
	err = stmt.QueryRow(username, noteid).Scan(&write)
	if err == sql.ErrNoRows {
		writePermission = false
		return writePermission
	}
	if err != nil {
		log.Panic(err)
	}

	//Check the permission
	if write == "t" {
		writePermission = true
		return writePermission
	}

	writePermission = false
	return writePermission

}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR HTTP TO CONTROLLER
// ------------------------------------------------------------------------------------------------------------------------------------------------

//Check if user is a owner of a note
func noteOwner(r *http.Request) bool {
	var owner string
	db := connectDatabase()
	defer db.Close()
	username := getUserName(r)
	noteID := mux.Vars(r)["id"]

	stmt, err := db.Prepare("SELECT note_owner FROM _note WHERE note_owner = $1 AND note_id = $2;")
	if err != nil {
		log.Panic(err)
	}
	err = stmt.QueryRow(username, noteID).Scan(&owner)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}

	return true

}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR HTTP TO CONTROLLER
// ------------------------------------------------------------------------------------------------------------------------------------------------

//get note based on note id and user permissions
func getPartOfNote(r *http.Request) (note Note) {
	//Check to see if user has read permission or is note owner
	if readPermissions(r) || noteOwner(r) {
		//Connect to database
		db := connectDatabase()
		defer db.Close()
		//Get username
		username := getUserName(r)
		noteid := mux.Vars(r)["id"]
		//prepare statment
		stmt, err := db.Prepare("SELECT _note.note_id, note_owner, title, body, date_created, read, write FROM _note JOIN _note_privileges ON _note.note_id = _note_privileges.note_id WHERE _note.note_id = $2 AND user_name = $1")
		if err != nil {
			log.Panic(err)
		}
		err = stmt.QueryRow(username, noteid).Scan(&note.NoteID, &note.NoteOwner, &note.NoteTitle, &note.NoteBody, &note.CreatedDate, &note.Read, &note.Write)
		if err == sql.ErrNoRows {
			return note
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return note
}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR HTTP TO CONTROLLER
// ------------------------------------------------------------------------------------------------------------------------------------------------

func getOwnedNote(r *http.Request) (note Note) {
	if noteOwner(r) {
		//Connect to database
		db := connectDatabase()
		defer db.Close()
		//get username
		username := getUserName(r)
		noteid := mux.Vars(r)["id"]
		//prepare statment
		stmt, err := db.Prepare("SELECT _note.note_id, note_owner, title, body, date_created FROM _note WHERE note_owner = $1 AND note_id = $2;")
		if err != nil {
			log.Panic(err)
			return note
		}
		err = stmt.QueryRow(username, noteid).Scan(&note.NoteID, &note.NoteOwner, &note.NoteTitle, &note.NoteBody, &note.CreatedDate)
		if err == sql.ErrNoRows {
			return note
		}
		if err != nil {
			log.Panic(err)
			return note
		}

	}
	return note
}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR HTTP TO CONTROLLER
// ------------------------------------------------------------------------------------------------------------------------------------------------

//Update note that you own
func updateOwnedNote(r *http.Request) (success bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//get username
	username := getUserName(r)
	//get the note id that we need to update
	noteid := mux.Vars(r)["id"]
	//get values from form
	title := r.FormValue("title")
	body := r.FormValue("body")
	fmt.Println("This is the title and body of update: " + title + " " + body) //for testing
	stmt, err := db.Prepare("UPDATE _note SET title = $1, body = $2  WHERE note_id = $3 AND note_owner = $4;")
	if err != nil {
		log.Panic(err)
	}

	result, err := stmt.Exec(title, body, noteid, username)
	if err != nil {
		log.Panic(err)
		success = false
		return success
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Panic(err)
		return false
	}
	if count > 0 {
		return true
	}
	return false

}

//Update note you are part of
func updatePartOfNote(r *http.Request) (success bool) {
	if userStillLoggedIn(r) {
		//Connect to database
		db := connectDatabase()
		defer db.Close()

		//get username
		//username := getUserName(r)
		//get note id
		noteID := mux.Vars(r)["id"]

		//Get value from form
		body := r.FormValue("body")
		fmt.Println("This is the update body: " + body) //This is for testing
		//prepare execution query
		stmt, err := db.Prepare("UPDATE _note SET body = $1 WHERE note_id = $2")
		result, err := stmt.Exec(body, noteID)
		if err != nil {
			log.Panic(err)
			success = false
			return success
		}
		//validate that update worked
		count, err := result.RowsAffected()
		if err != nil {
			log.Panic(err)
			success = false
			return success
		}
		if count > 0 {
			success = true
			return success
		}
	}
	success = false
	return success
}

func getOwnedNotesSQL(username string) (notes []Note) {
	//Connect to Database
	db := connectDatabase()
	defer db.Close()
	var note Note
	//Prepare Statment
	stmt, err := db.Prepare("SELECT _note.note_id, title, body, date_created FROM _note WHERE note_owner=$1")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query(username)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody, &note.CreatedDate)
		notes = append(notes, note)
	}

	return notes
}

func getPartOfNotesSQL(username string) (notes []Note) {
	db := connectDatabase()
	defer db.Close()

	var note Note
	//prepare statement
	stmt, err := db.Prepare(`
	SELECT _note.note_id, title, body, date_created, note_owner FROM _note_privileges
	JOIN _note
	ON _note.note_id = _note_privileges.note_id
	WHERE _note_privileges.user_name = $1;
	`)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query(username)
	//scan each row of the query and add it to the notes slice
	for rows.Next() {
		rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody, &note.CreatedDate, &note.NoteOwner)
		log.Println("Notes part of:" + note.NoteTitle) //for testing
		notes = append(notes, note)
	}
	return notes
}
