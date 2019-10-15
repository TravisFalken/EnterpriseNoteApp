package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//=========================Checks if user login details are correct=========================================
func login(w http.ResponseWriter, r *http.Request) {
	var loginUser User
	req, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(req, &loginUser)
	if userNameExists(loginUser.UserName) {

		if validatePass(loginUser.Password) {

			//create new session id
			sessionid := newSessionid()
			//Add session id to the user in the database
			addSessionToUser(loginUser, sessionid)
			//set the clients session cookie
			sessionCookie := &http.Cookie{
				Name:  "session",
				Value: sessionid,
			}
			http.SetCookie(w, sessionCookie)
			//Create username cookie
			usernameCookie := &http.Cookie{
				Name:  "username",
				Value: loginUser.UserName,
			}
			//Set a cookie to username
			http.SetCookie(w, usernameCookie)
			fmt.Fprintf(w, "Successfully logged in")
		} else {
			fmt.Fprintf(w, "Login not successfull")
		}
	} else {
		fmt.Fprintf(w, "Login not successfull")
	}

}

//====================ADD USER=====================================
func addUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered addUser()") // For testing

	var newUser User
	//Get user information out of body of HTTP
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newUser)

	fmt.Fprintf(w, addUserSQL(newUser))

}

//==========================ADD NOTE=============================================================
func addNote(w http.ResponseWriter, r *http.Request) {
	//Check if user is logged in
	if userStillLoggedIn(r) {
		var newNote Note
		var noteTime = time.Now()
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			log.Fatal(err)
		}
		username := usernameCookie.Value

		//Get the body and put it into a a note struct
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &newNote)

		//set the created date of the note
		newNote.CreatedDate = noteTime.String()
		newNote.NoteOwner = username

		fmt.Fprintf(w, addNoteSQL(newNote))

	} else {
		fmt.Fprintf(w, "You are not logged in!")
	}

}

/*===============================LIST ALL NOTES BELONGING TO USER==================================*/

func listNotes(w http.ResponseWriter, r *http.Request) {
	//Check if user is still online
	if userStillLoggedIn(r) {

		usernameCookie, err := r.Cookie("username")
		if err != nil {
			log.Fatal(err)
		}
		username := usernameCookie.Value
		notes := listAllNotesSQL(username)

		fmt.Println(notes)
		//just sending it straight to frontend for testing
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
		//User is not logged in
	} else {
		fmt.Fprintf(w, "Not Logged in!")
	}

}

//============================================LIST ALL REGISTERED USERS=====================
//This function creates a list of all registerd users
func listAllUses(loggedInUser string) []string {

	return listAllUsersSQL(loggedInUser)
}

//=========================USER LOGOUT======================================================
func logout(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		deleteSesion(w, r)
		fmt.Fprintf(w, "Successfully logged out")
	} else {
		fmt.Fprintf(w, "Already Logged out")
	}

}

//==============DELETE A NOTE==============================
func deleteNote(w http.ResponseWriter, r *http.Request) {
	//Check if user is still logged in
	if userStillLoggedIn(r) {
		noteid := mux.Vars(r)["id"]
		username := getUserName(r)
		//deletes a note and returns true if deleted
		if deleteSpecificNote(noteid, username) {
			fmt.Fprintf(w, "Successfully Deleted")
		} else {
			fmt.Fprintf(w, "Not Successful")
		}
	} else {
		fmt.Fprintf(w, "You cannot delete note because you are not logged in")
	}
}

//==============SEARCH FOR A NOTE WITH PARTIAL TEXT SQL==============================
func searchNotePartial(w http.ResponseWriter, r *http.Request) {
	//Check if user is still online
	if userStillLoggedIn(r) {

		bodyText := mux.Vars(r)["id"]
		username := getUserName(r)
		notes := partialTextSearchSQL(bodyText, username)

		fmt.Println(notes)
		//just sending it straight to frontend for testing
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
		//User is not logged in
	} else {
		fmt.Fprintf(w, "Not Logged in!")
	}

}
