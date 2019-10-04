package main

import (
	"database/sql"
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
func validatePass(password string) bool {
	var pass string
	//Connect to db
	db := connectDatabase()
	defer db.Close()

	//prepare statement to check for password
	stmt, err := db.Prepare("SELECT password FROM _user WHERE password = $1;")
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(password).Scan(&pass)
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

//====================ADD USER=====================================
func addUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered addUser()") // For testing

	var newUser User
	//Get user information out of body of HTTP
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newUser)

	fmt.Fprintf(w, addUserSQL(newUser))

	//Create connection to server
	//Connect to db
	// db := connectDatabase()
	// defer db.Close()

	// if !userNameExists(newUser.UserName) {
	// 	//Prepare insert to stop SQL injections
	// 	log.Println("Entered add user if statement")
	// 	stmt, err := db.Prepare("INSERT INTO _user VALUES($1,$2,$3,$4,$5);")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	_, err = stmt.Exec(newUser.UserName, newUser.Password, newUser.Email, newUser.GivenName, newUser.FamilyName)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Fprintf(w, "Added user")
	// } else {
	// 	fmt.Fprintf(w, "Username already exists")
	// }

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

		//Connect to db
		// db := connectDatabase()
		// defer db.Close()

		// //Prepare insert to stop SQL injections
		// stmt, err := db.Prepare("INSERT INTO _note (title, body, date_created, note_owner) VALUES($1,$2,$3,$4);")
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// _, err = stmt.Exec(newNote.NoteTitle, newNote.NoteBody, newNote.CreatedDate, newNote.NoteOwner)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// fmt.Fprintf(w, "New Note Added")
		// //User is not logged in
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
		//var notes []Note
		//Connect to db
		// db := connectDatabase()
		// defer db.Close()
		// stmt, err := db.Prepare("SELECT _note.note_id, _note.note_owner, _note.title, _note.body, _note.date_created FROM _note LEFT OUTER JOIN _note_privileges ON (_note.note_id = _note_privileges.note_id) WHERE _note.note_owner = $1 OR _note_privileges.user_name = $1;")
		// var note Note
		// rows, err := stmt.Query(username)
		// for rows.Next() {
		// 	err = rows.Scan(&note.NoteID, &note.NoteOwner, &note.NoteTitle, &note.NoteBody, &note.CreatedDate)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	notes = append(notes, note)

		// }
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
	// var users []string
	// var username string
	// //Connect to db
	// db := connectDatabase()
	// defer db.Close()

	// //Send query to the db
	// rows, err := db.Query("SELECT user_name FROM _user;")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for rows.Next() {
	// 	err = rows.Scan(&username)
	// 	//Make sure we dont add logged in user to recommended
	// 	if loggedInUser != username {
	// 		users = append(users, username)
	// 	}
	// }

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
		//deletes a note and returns true if deleted
		if deleteSpecificNote(r) {
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
		//Connect to db
		// db := connectDatabase()
		// defer db.Close()

		// bodyText += ":*" //for testing
		// stmt, err := db.Prepare("SELECT _note.note_id, _note.note_owner, _note.title, _note.body, _note.date_created FROM _note LEFT OUTER JOIN _note_privileges ON (_note.note_id = _note_privileges.note_id) WHERE body ~ $2 AND _note.note_owner = $1 OR _note_privileges.user_name = $1;")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// var note Note
		// rows, err := stmt.Query(username, bodyText)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// for rows.Next() {
		// 	err = rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody, &note.CreatedDate, &note.NoteOwner)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	notes = append(notes, note)

		// }
		fmt.Println(notes)
		//just sending it straight to frontend for testing
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
		//User is not logged in
	} else {
		fmt.Fprintf(w, "Not Logged in!")
	}

}
