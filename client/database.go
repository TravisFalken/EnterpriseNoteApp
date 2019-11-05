package main

import (
	"database/sql"
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

func readPermissionsSQL(username string, noteid string, read string) bool {
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
		return false
	}
	if read == "t" {

		return true
	}

	return false

}

func checkWritePermissionsSQL(username string, noteid string, write string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//prepare statement
	stmt, err := db.Prepare("SELECT write FROM _note_privileges WHERE user_name = $1 AND note_id = $2;")
	if err != nil {
		log.Panic(err)
	}
	err = stmt.QueryRow(username, noteid).Scan(&write)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Panic(err)
	}

	//Check the permission
	if write == "t" {
		return true
	}

	return false
}

func noteOwnerSQL(username string, noteid string, owner string) bool {
	db := connectDatabase()
	defer db.Close()

	stmt, err := db.Prepare("SELECT note_owner FROM _note WHERE note_owner = $1 AND note_id = $2;")
	if err != nil {
		log.Panic(err)
	}
	err = stmt.QueryRow(username, noteid).Scan(&owner)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}

	return true
}

//get all users that are not already part of note
func getAvaliableUsers(r *http.Request) (users []string) {
	//connect to database
	db := connectDatabase()
	defer db.Close()
	var user string
	username := getUserName(r)
	noteid := mux.Vars(r)["id"]

	//prepare statement
	stmt, err := db.Prepare("SELECT _user.user_name FROM _user WHERE  _user.user_name NOT IN (SELECT user_name FROM _note_privileges WHERE _note_privileges.note_id = $1)")
	//"SELECT user_name FROM _user JOIN _note_privilges ON _user.user_name = _note_privileges.user_name WHERE _note_privileges.note_id = $1 AND _note_privilegs.user_name = NULL;")
	if err != nil {
		log.Panic(err)
		return users
	}

	rows, err := stmt.Query(noteid)
	if err != nil {
		log.Panic(err)
		return users
	}
	for rows.Next() {
		rows.Scan(&user)
		if user != username {
			users = append(users, user)
		}
	}
	return users
}

//get all the privilges from a note
func getNotePrivileges(noteid string) (privileges []privlige) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	var newPrivilege privlige

	//prepare statment
	stmt, err := db.Prepare("SELECT user_name, read, write FROM _note_privileges WHERE note_id = $1;")
	if err != nil {
		log.Panic(err)
		return privileges
	}

	rows, err := stmt.Query(noteid)
	if err != nil {
		log.Panic(err)
		return privileges
	}

	for rows.Next() {
		rows.Scan(&newPrivilege.Username, &newPrivilege.Read, &newPrivilege.Write)
		privileges = append(privileges, newPrivilege)
	}
	return privileges
}

//Remove a privilege from a note
func removePrivilege(noteid string, username string) bool {
	//Connect to Database
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare("DELETE FROM _note_privileges WHERE note_id = $1 AND user_name = $2;")
	if err != nil {
		log.Panic(err)
		return false
	}

	deleted, _ := stmt.Exec(noteid, username)
	//Validate if any row has been deleted
	rowsAffected, _ := deleted.RowsAffected()
	if rowsAffected > 0 {
		return true
	}
	return false
}

//Update exisiting privilege for a note
func updatePrivilege(noteid string, username string, write string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("UPDATE _note_privileges SET write = $1 WHERE note_id = $2 AND user_name = $3;")
	if err != nil {
		log.Panic(err)
		return false
	}

	result, err := stmt.Exec(write, noteid, username)
	if err != nil {
		log.Panic(err)
		return false
	}
	//Validate that the privilege has been updated
	count, _ := result.RowsAffected()
	if count > 0 {
		return true
	}
	return false
}
