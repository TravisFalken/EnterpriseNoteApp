package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

//==================gets the user from database with the session id================================
func getUser(sessionid string) (user *User, userFound bool) {

	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//Prepare query for getting user with the session id
	stmt, err := db.Prepare("SELECT username,email,givenname FROM _user WHERE sessionid = $1;")
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
		err = rows.Scan(&user)
		if err != nil {
			log.Fatal(err)
		}

	}
	//Checks if no user has been returned does not match session id
	if user == nil {
		userFound = false
		return user, userFound
	}
	userFound = true
	return user, userFound

}

//==============Creates a new Session id================================================
func newSessionid() string {
	id, _ := uuid.NewV4()
	return id.String()
}

//=======================Add session id to user======================================================
func addSessionToUser(user User, sessionID string) bool {

	//Connect to db
	db := connectDatabase()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE _user SET sessionid =$1 WHERE username=$2;")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(sessionID, user.UserName)
	return true
}

//======================Check if user is still logged in========================================================
func userStillLoggedIn(req *http.Request) bool {
	var username string
	sessionCookie, err := req.Cookie("session")
	log.Println(sessionCookie.Value) //For testing
	if err != nil {
		return false
	}

	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//Prepare statement to stop sql injection
	stmt, err := db.Prepare("SELECT username FROM _user WHERE sessionid=$1;")
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(sessionCookie.Value).Scan(&username)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	return true
}

//Delete session id from  database and users computer
func deleteSesion(w http.ResponseWriter, r *http.Request) bool {
	sessionid, err := r.Cookie("session")
	if err != nil {
		return true
	}
	sessionid = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, sessionid)
	return true
}
