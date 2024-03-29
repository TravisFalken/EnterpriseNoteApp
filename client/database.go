package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func connectDatabase() (db *sql.DB) {
	// Connection string needed for docker
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
	getUserName, err := db.Prepare(`
		Select user_name 
		FROM _user 
		WHERE user_name = $1;`)
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
	stmt, err := db.Prepare(`
		SELECT password 
		FROM _user 
		WHERE password = $1 AND user_name = $2;`)
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
func addUserSQL(newUser User) bool {

	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Validate that username does not already exist
	if !userNameExists(newUser.UserName) {
		//Prepare insert to stop SQL injections
		stmt, err := db.Prepare(`
			INSERT INTO _user(user_name, password, email, given_name, family_name)
			VALUES($1,$2,$3,$4,$5);`)
		if err != nil {
			log.Panic(err)
			return false
		}

		_, err = stmt.Exec(newUser.UserName, newUser.Password, newUser.Email, newUser.GivenName, newUser.FamilyName)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}

	return false

}

//===================Add Note=============================
//

func addNoteSQL(newNote Note) (noteid string) {

	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//Prepare insert to stop SQL injections
	stmt, err := db.Prepare(`
		INSERT INTO _note (title, body, date_created, note_owner) 
		VALUES($1,$2,$3,$4) 
		RETURNING note_id;`)
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(newNote.NoteTitle, newNote.NoteBody, newNote.CreatedDate, newNote.NoteOwner).Scan(&noteid)
	if err != nil {
		log.Fatal(err)
	}

	return noteid

}

//===================List All Notes=============================
//

func listAllNotesSQL(username string) []Note {

	var notes []Note
	//Connect to db
	db := connectDatabase()
	defer db.Close()
	stmt, err := db.Prepare(`
		SELECT _note.note_id, _note.note_owner, _note.title, _note.body, _note.date_created 
		FROM _note 
		LEFT OUTER JOIN _note_privileges ON (_note.note_id = _note_privileges.note_id) 
		WHERE _note.note_owner = $1 OR _note_privileges.user_name = $1;`)
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
	rows, err := db.Query(`
		SELECT user_name 
		FROM _user;`)
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

	bodyText += ":*"
	stmt, err := db.Prepare(`
		SELECT _note.note_id, _note.title, _note.body, _note.date_created, _note.note_owner 
		FROM _note 
		LEFT OUTER JOIN _note_privileges ON (_note.note_id = _note_privileges.note_id) 
		WHERE body ~ $2 AND _note.note_owner = $1 OR _note_privileges.user_name = $1;`)
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
func partialSeachOwnedTitleSQL(searchText string, username string) (ownedNotes []Note) {

	//Connect to database

	db := connectDatabase()
	defer db.Close()
	searchText += ":*"
	stmt, err := db.Prepare(`
		SELECT _note.note_id, _note.title, _note.body, _note.date_created, _note.note_owner 
		FROM _note 
		WHERE _note.title ~ $2 AND _note.note_owner = $1;`)
	if err != nil {
		log.Fatal(err)
	}
	var note Note
	rows, err := stmt.Query(username, searchText)
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody, &note.CreatedDate, &note.NoteOwner)

		if err != nil {
			log.Panic(err)
		}
		ownedNotes = append(ownedNotes, note)
	}
	return ownedNotes
}

//=============Partial search notes you are appart of by their title=================
func partialSearchPartOfTitleSQL(titleText string, username string) (partOfNotes []Note) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	titleText += ":*"
	stmt, err := db.Prepare(`
		SELECT _note.note_id, _note.title, _note.body, _note.date_created, _note.note_owner 
		FROM _note_privileges 
		JOIN _note ON _note_privileges.note_id = _note.note_id 
		WHERE _note.title ~* $2 AND _note_privileges.user_name = $1`)
	if err != nil {
		log.Fatal(err)
	}
	var note Note
	rows, err := stmt.Query(username, titleText)
	if err != nil {
		log.Panic(err)
		return partOfNotes
	}
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody, &note.CreatedDate, &note.NoteOwner)
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

	stmt, err := db.Prepare(`
		DELETE FROM _note 
		WHERE note_owner=$1 AND note_id=$2;`)
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

	stmt, err := db.Prepare(`
		DELETE FROM _note 
		WHERE note_owner=$1;`)
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

	stmt, err := db.Prepare(`
		DELETE FROM _user 
		WHERE user_name=$1;`)
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

//===================Get Owned Notes=============================

func getOwnedNotesSQL(username string) (notes []Note) {
	//Connect to Database
	db := connectDatabase()
	defer db.Close()
	var note Note
	//Prepare Statment // needs owner for testing
	stmt, err := db.Prepare(`
		SELECT _note.note_id, title, body, date_created, note_owner 
		FROM _note 
		WHERE note_owner=$1`)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query(username)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody, &note.CreatedDate, &note.NoteOwner)
		notes = append(notes, note)
	}

	return notes
}

//===================Get Part Of Notes=============================

func getPartOfNotesSQL(username string) (notes []Note) {
	db := connectDatabase()
	defer db.Close()

	var note Note
	//prepare statement
	stmt, err := db.Prepare(`
		SELECT _note.note_id, title, body, date_created, note_owner FROM _note_privileges
		JOIN _note
		ON _note.note_id = _note_privileges.note_id
		WHERE _note_privileges.user_name = $1;`)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query(username)
	//scan each row of the query and add it to the notes slice
	for rows.Next() {
		rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody, &note.CreatedDate, &note.NoteOwner)
		notes = append(notes, note)
	}
	return notes
}

//===================Get Single Owned Note=============================

func getOwnedNoteSQL(noteid string, username string) (ownedNote Note) {

	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//prepare statment
	stmt, err := db.Prepare(`
		SELECT _note.note_id, note_owner, title, body, date_created 
		FROM _note 
		WHERE note_owner = $1 AND note_id = $2;`)
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
	return note
}

//===================Get Single Part Of Note=============================

func getPartOfNoteSQL(noteid string, username string) (note Note) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//prepare statment
	stmt, err := db.Prepare(`
		SELECT _note.note_id, note_owner, title, body, date_created, read, write 
		FROM _note 
		JOIN _note_privileges ON _note.note_id = _note_privileges.note_id 
		WHERE _note.note_id = $2 AND user_name = $1`)
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
	return note
}

//===================Get read permissions=============================

func readPermissionsSQL(username string, noteid string, read string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	stmt, err := db.Prepare(`
		SELECT read 
		FROM _note_privileges 
		WHERE user_name = $1 AND note_id = $2`)
	//if no rows were returned
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Panic(err)
	}
	err = stmt.QueryRow(username, noteid).Scan(&read)
	//if no rows were returned
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		return false
	}
	if read == "t" {

		return true
	}

	return false

}

//===================Check Permissions============================

func checkWritePermissionsSQL(username string, noteid string, write string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//prepare statement
	stmt, err := db.Prepare(`
		SELECT write 
		FROM _note_privileges 
		WHERE user_name = $1 AND note_id = $2;`)
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

	if write == "t" {
		return true
	}

	return false
}

//=================== WRONG =============================
// ----------------------------------------------------------------------------------------------------------------------
// Travis this code doesnt make sense. you are selecting note owner where note owner equals username you are passing to query. therefore getting a return of something that you already have
// ---------------------------------------------------------------------------
func noteOwnerSQL(username string, noteid string, owner string) bool {
	db := connectDatabase()
	defer db.Close()

	stmt, err := db.Prepare(`
		SELECT note_owner 
		FROM _note 
		WHERE note_owner = $1 AND note_id = $2;`)
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
func getAvaliableUsersSQL(username string, noteid string) (users []string) {
	//connect to database
	db := connectDatabase()
	defer db.Close()
	var user string

	//prepare statement
	stmt, err := db.Prepare(`
		SELECT _user.user_name 
		FROM _user 
		WHERE  _user.user_name NOT IN (
			SELECT user_name 
			FROM _note_privileges 
			WHERE _note_privileges.note_id = $1)`)
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
func getNotePrivileges(noteid string) (privileges []Privlige) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	var newPrivilege Privlige

	//prepare statment
	stmt, err := db.Prepare(`
		SELECT user_name, read, write 
		FROM _note_privileges 
		WHERE note_id = $1;`)
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
	stmt, err := db.Prepare(`
		DELETE FROM _note_privileges 
		WHERE note_id = $1 AND user_name = $2;`)
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
	stmt, err := db.Prepare(`
		UPDATE _note_privileges 
		SET write = $1 
		WHERE note_id = $2 AND user_name = $3;`)
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

//Sql For updating ownded Note
func updateOwnedNoteSQL(title string, body string, noteid string, noteOwner string) (success bool) {
	//Connect to Database
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare(`
		UPDATE _note 
		SET title = $1, body = $2  
		WHERE note_id = $3 AND note_owner = $4;`)
	if err != nil {
		log.Panic(err)
	}

	result, err := stmt.Exec(title, body, noteid, noteOwner)
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

// -------------------------------------------------------------------------------------------------
// ERROR HERE. This code is exactly the same as any update code. i understand that it is after an part of check in controller however it would be a better idea to have a generic update note query that is called to remove double ups
// -------------------------------------------------------------------------------------------------

//Sql for updating a note user is part of
func updatePartOfNoteSQL(noteID string, body string) (success bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	stmt, err := db.Prepare(`
		UPDATE _note 
		SET body = $1 
		WHERE note_id = $2`)
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
	success = false
	return success
}

//Add Permissions SQL
func addPermissionSQL(noteid string, user string, read string, write string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//Prepare statment
	stmt, err := db.Prepare(`
		INSERT INTO _note_privileges(note_id,user_name, read, write) 
			VALUES($1, $2, $3, $4);`)
	if err != nil {
		log.Panic(err)
		return false
	}
	_, err = stmt.Exec(noteid, user, read, write)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

////////////////SAVED PERMISSIONS SECTION///////////////////////

//Create a new group
func createNewGroup(groupName string, groupOwner string, read string, write string) (groupID string) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare(`
		INSERT INTO _group(group_title, read, write, group_owner) 
		VALUES($1,$2,$3,$4) 
		RETURNING group_id;`)
	if err != nil {
		log.Panic(err)
	}
	err = stmt.QueryRow(groupName, read, write, groupOwner).Scan(&groupID)
	if err != nil {
		log.Fatal(err)
	}
	return groupID
}

//get one group based on group id
func getGroup(groupid string) (group Group) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare(`
		SELECT group_id, group_title, read, write 
		FROM _group 
		WHERE group_id = $1;`)
	if err != nil {
		log.Panic(err)
	}

	err = stmt.QueryRow(groupid).Scan(&group.GroupID, &group.GroupTitle, &group.GroupRead, &group.GroupWrite)
	if err == sql.ErrNoRows {
		log.Panic(err)
		return group
	}
	if err != nil {
		log.Fatal(err)
	}

	return group
}

//Get all of the groups
func getAllGroups(username string) (groups []Group) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	var group Group
	//Prepare statement
	stmt, err := db.Prepare(`
		SELECT group_id, group_title, read, write 
		FROM _group 
		WHERE group_owner = $1;`)
	if err != nil {
		log.Panic(err)
	}
	rows, err := stmt.Query(username)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&group.GroupID, &group.GroupTitle, &group.GroupRead, &group.GroupWrite)
		groups = append(groups, group)
	}

	return groups
}

//Save a user to a group
func saveGroupUserSQL(groupID string, username string) bool {
	//Connect to DB
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare(`
		INSERT INTO _group_user(group_id, user_name) 
		VALUES($1,$2);`)
	if err != nil {
		log.Panic(err)
	}

	_, err = stmt.Exec(groupID, username)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

//Validate if user is group owner
func validateGroupOwner(username string, groupid string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	var groupUsername string

	//Prepare statement
	stmt, err := db.Prepare(`
		SELECT group_owner 
		FROM _group 
		WHERE group_id = $1;`)
	if err != nil {
		log.Panic(err)
	}

	err = stmt.QueryRow(groupid).Scan(&groupUsername)
	if err == sql.ErrNoRows {
		log.Panic(err)
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	//Validate that username equals group username
	if username != groupUsername {
		return false
	}

	return true
}

//get all saved users for a group
func getGroupUsers(groupid string) (users []string) {

	var user string
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare(`
		SELECT user_name 
		FROM _group_user 
		WHERE group_id = $1;`)
	if err != nil {
		log.Panic(err)
	}

	rows, err := stmt.Query(groupid)
	if err != nil {
		log.Fatal(err)
	}
	//run through all the rows of the query
	for rows.Next() {
		rows.Scan(&user)
		users = append(users, user)
	}

	return users
}

//get group privileges
func getGroupPrivileges(groupid string) (read string, write string) {
	//Connect Database
	db := connectDatabase()
	defer db.Close()

	stmt, err := db.Prepare(`
		SELECT read, write 
		FROM _group 
		WHERE group_id = $1;`)
	if err != nil {
		log.Panic(err)
	}
	stmt.QueryRow(groupid).Scan(&read, &write)
	return read, write
}

//edit group privileges
func editGroupPrivileges(groupid string, write string, read string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare(`
		UPDATE _group 
		SET write = $1, read = $2 
		WHERE group_id = $3;`)
	if err != nil {
		log.Panic(err)
	}
	result, err := stmt.Exec(write, read, groupid)
	if err != nil {
		log.Panic(err)
		return false
	}

	//Validate that row updated
	count, err := result.RowsAffected()
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Panic(err)
		return false
	}
	if count == 0 {
		return false
	}

	return true
}

//Remove user from group
func removeGroupUser(groupid string, user string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare(`
		DELETE FROM _group_user 
		WHERE group_id = $1 AND user_name = $2;`)
	if err != nil {
		log.Panic(err)
	}

	result, err := stmt.Exec(groupid, user)
	if err != nil {
		log.Fatal(err)
	}

	//validate that user has been deleted
	count, err := result.RowsAffected()
	if err == sql.ErrNoRows {
		log.Panic(err)
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		return false
	}
	return true
}

//Remove group from database
func removeGroup(groupid string) bool {
	//Connect database
	db := connectDatabase()

	//Prepare statement
	stmt, err := db.Prepare(`
		DELETE FROM _group 
		WHERE group_id = $1;`)
	if err != nil {
		log.Panic(err)
	}
	result, err := stmt.Exec(groupid)
	if err != nil {
		log.Fatal(err)
	}
	//Validate that note has been deleted
	count, err := result.RowsAffected()
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Panic(err)
	}
	if count > 0 {
		return true
	}
	return false
}

//edit a group and save it in the database
func editGroup(groupid string, write string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare(`
		UPDATE _group 
		SET write = $1 
		WHERE group_id = $2;`)
	if err != nil {
		log.Panic(err)
	}

	result, err := stmt.Exec(write, groupid)
	if err != nil {
		log.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err == sql.ErrNoRows {
		log.Panic(err)
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return true
	}
	return false
}

//Get all of the users that are not part of the group
func getAvaliableGroupUsers(groupid string, username string) (users []string) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	var user string
	//prepare statement
	stmt, err := db.Prepare(`
		SELECT _user.user_name 
		FROM _user 
		WHERE  _user.user_name 
		NOT IN (
			SELECT user_name 
			FROM _group_user 
			WHERE _group_user.group_id = $1)`)
	if err != nil {
		log.Panic(err)
	}

	rows, err := stmt.Query(groupid)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&user)
		if err != nil {
			log.Panic(err)
		}
		//Make sure note owner is note part of users
		if user != username {
			users = append(users, user)
		}
	}

	return users
}

//=======================Add session id to user======================================================
func addSessionToUser(user User, sessionID string) bool {

	//Connect to db
	db := connectDatabase()
	defer db.Close()

	stmt, err := db.Prepare(`
		UPDATE _user 
		SET session_id=$1 
		WHERE user_name=$2;`)
	if err != nil {
		log.Fatal(err)
	}
	result, err := stmt.Exec(sessionID, user.UserName)
	if err != nil {
		log.Fatal(err)
	}

	//Validate that session has been added to user
	count, err := result.RowsAffected()
	//Session was not added to user
	if err == sql.ErrNoRows {
		log.Panic(err)
		return false
	}
	if err != nil {
		log.Panic(err)
	}
	//if count greater than 0 session has been added to user
	if count > 0 {
		return true
	}
	return false
}

//==================gets the user from database with the session id================================
func getUser(sessionid string) (user User) {

	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//Prepare query for getting user with the session id
	stmt, err := db.Prepare(`
		SELECT user_name,given_name,family_name, email, password, session_id  
		FROM _user WHERE session_id = $1;`)
	if err != nil {
		log.Fatal(err)
	}
	//Query the database to get a user that matches the session
	rows, err := stmt.Query(sessionid)
	if err != nil {
		log.Fatal(err)
	}
	//run through the row
	for rows.Next() {
		err = rows.Scan(&user.UserName, &user.GivenName, &user.FamilyName, &user.Email, &user.Password, &user.SessionID)
		if err != nil {
			log.Fatal(err)
		}

	}

	return user

}
