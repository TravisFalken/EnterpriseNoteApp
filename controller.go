package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//This method runs the first time the user trys to access the webapp
//Looking for a better function name XD
func index(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//============This is the home page===============================
func home(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		sessionid := getSession(r)
		user := getUser(sessionid)
		log.Println(user)
		err := tpl.ExecuteTemplate(w, "home.gohtml", user)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//===========================THE SIGNUP Page============================================
func signUp(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	} else {
		tpl.ExecuteTemplate(w, "signup.gohtml", nil)
	}

}

//=====================THE CREATE NOTE PAGE===========================================
func createNote(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		tpl.ExecuteTemplate(w, "createNote.gohtml", nil)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//=========================Checks if user login details are correct=========================================
func login(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {

		var loginUser User
		loginUser.UserName = r.FormValue("username")
		loginUser.Password = r.FormValue("password")
		if userNameExists(loginUser.UserName) {

			if validatePass(loginUser.Password, loginUser.UserName) {

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
				//Send the home page to user
				log.Println("successfully logged in") // for testing
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				http.Error(w, "username and/or password does not match", http.StatusForbidden)
				return
			}
		} else {
			http.Error(w, "username and/or password does not match", http.StatusForbidden)
			return
		}
	}
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//====================ADD USER=====================================
func addUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Entered addUser()") // For testing // old add user

	// var newUser User
	// //Get user information out of body of HTTP
	// reqBody, _ := ioutil.ReadAll(r.Body)
	//json.Unmarshal(reqBody, &newUser)

	// new add user
	fmt.Println("Entered addUser()") // For testing

	var newUser User
	newUser.UserName = r.FormValue("user_name")
	newUser.GivenName = r.FormValue("given_name")
	newUser.FamilyName = r.FormValue("family_name")
	newUser.Email = r.FormValue("email")
	newUser.Password = r.FormValue("password")
	if !userNameExists(newUser.UserName) {
		fmt.Fprintf(w, addUserSQL(newUser))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		fmt.Fprintf(w, "Username already exists")
	}

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
		newNote.NoteTitle = r.FormValue("title")
		newNote.NoteBody = r.FormValue("body")

		newNote.CreatedDate = noteTime.Format("2006-01-02")
		log.Println(newNote.CreatedDate) // For testing
		newNote.NoteOwner = username

		fmt.Fprintf(w, addNoteSQL(newNote))

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		//fmt.Fprintf(w, "Successfully logged out")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		//fmt.Fprintf(w, "Already logged out")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

//==============List All Notes that user owners and is part of==================================
func allNotes(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		var username = getUserName(r)
		var notes Notes
		notes.OwnedNotes = getOwndedNotes(username)
		notes.PartOfNotes = getPartOfNotes(username)
		log.Println(notes) //For testing
		tpl.ExecuteTemplate(w, "listNotes.gohtml", notes)

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

//==========GET NOTES That you are only appart of====================
//Still need to do
func getPartOfNotes(username string) (notes []Note) {
	return notes
}
