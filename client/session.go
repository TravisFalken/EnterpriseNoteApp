package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

//==============Creates a new Session id================================================
func newSessionid() string {
	id, _ := uuid.NewV4()
	return id.String()
}

//======================Check if user is still logged in========================================================
func userStillLoggedIn(req *http.Request) bool {
	sessionid := getSession(req)
	user := getUser(sessionid)
	log.Println("Database Sessionid And Client ID: " + user.SessionID + "|||" + sessionid) // for testing
	if user.SessionID == sessionid {
		return true
	}

	return false
	/*
		//Connect to db
		db := connectDatabase()
		defer db.Close()

		//Prepare statement to stop sql injection
		stmt, err := db.Prepare("SELECT user_name FROM _user WHERE session_id=$1;")
		if err != nil {
			log.Fatal(err)
		}

		err = stmt.QueryRow(sessionid).Scan(&username)
		if err == sql.ErrNoRows {
			return false
		}
		if err != nil {
			log.Fatal(err)
		}
		return true
	*/
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
	sessionID := getSession(req)

	//get the users
	user := getUser(sessionID)
	return user.UserName
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
