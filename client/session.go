package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

//==================gets the user from database with the session id================================
func getUser(sessionid string) (user User) {

	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//Prepare query for getting user with the session id
	stmt, err := db.Prepare("SELECT user_name,given_name,family_name FROM _user WHERE session_id = $1;")
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
		err = rows.Scan(&user.UserName, &user.GivenName, &user.FamilyName)
		if err != nil {
			log.Fatal(err)
		}

	}

	return user

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

	stmt, err := db.Prepare("UPDATE _user SET session_id=$1 WHERE user_name=$2;")
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
	if err != nil {
		return false
	}

	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//Prepare statement to stop sql injection
	stmt, err := db.Prepare("SELECT user_name FROM _user WHERE session_id=$1;")
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

//Get username based on session id
func getUserName(req *http.Request) (username string) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	sessionID := getSession(req)

	//Prepare Query
	stmt, err := db.Prepare("SELECT user_name FROM _user WHERE session_id=$1;")
	if err != nil {
		log.Fatal(err)
	}

	//Query DB
	rows, err := stmt.Query(sessionID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			log.Fatal(err)
		}
	}
	return username
}

//===================GET sessionid as string==================================
func getSession(r *http.Request) (sessionid string) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		sessionid = " "
		return sessionid
	}
	sessionid = sessionCookie.Value
	return sessionid

}