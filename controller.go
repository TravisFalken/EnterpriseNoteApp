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
		user := getUserName(r)
		users := listAllUses(user)
		tpl.ExecuteTemplate(w, "createNote.gohtml", users)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//====================EDIT NOTE PAGE======================================
func editNote(w http.ResponseWriter, r *http.Request) {
	//Make sure user is still logged in
	if userStillLoggedIn(r) {
		//validate that user can actually read the note
		if readPermissions(r) || noteOwner(r) {
			var note Note
			//get note if user is part of
			if readPermissions(r) {
				note = getPartOfNote(r)
				//get note if user us owner
			} else {
				note = getOwnedNote(r)
			}
			//get the note
			log.Println("Edit Note: " + note.NoteTitle + note.Write + "Note Owner: " + note.NoteOwner) // For testing
			tpl.ExecuteTemplate(w, "editNote.gohtml", note)
		} else {
			http.Redirect(w, r, "/listNotes", http.StatusSeeOther)
		}

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

	//Not Sure
	/*
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
	*/
	fmt.Println("Entered addUser()") // For testing

	var newUser User
	//Get user information out of body of HTTP
	/*
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &newUser)
	*/

	newUser.UserName = r.FormValue("user_name")
	newUser.GivenName = r.FormValue("given_name")
	newUser.FamilyName = r.FormValue("family_name")
	newUser.Email = r.FormValue("email")
	newUser.Password = r.FormValue("password")
	//Create connection to server
	//Connect to db
	db := connectDatabase()
	defer db.Close()

	if !userNameExists(newUser.UserName) {
		//Prepare insert to stop SQL injections
		log.Println("Entered add user if statement")
		stmt, err := db.Prepare("INSERT INTO _user(user_name, password, email, given_name, family_name) VALUES($1,$2,$3,$4,$5);")
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(newUser.UserName, newUser.Password, newUser.Email, newUser.GivenName, newUser.FamilyName)
		if err != nil {
			log.Fatal(err)
		}
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

		//Connect to db
		db := connectDatabase()
		defer db.Close()

		//Prepare insert to stop SQL injections
		var noteId string
		stmt, err := db.Prepare("INSERT INTO _note (title, body, date_created, note_owner) VALUES($1,$2,$3,$4) RETURNING note_id;")
		if err != nil {
			log.Fatal(err)
		}

		err = stmt.QueryRow(newNote.NoteTitle, newNote.NoteBody, newNote.CreatedDate, newNote.NoteOwner).Scan(&noteId)
		if err != nil {
			log.Fatal(err)
		}
		//Get users attached to note and add them to the database
		var read string
		var write string

		stmt, err = db.Prepare("INSERT INTO _note_privileges(note_id, user_name, read, write) VALUES($1,$2,$3,$4);")
		if err != nil {
			log.Fatal(err)
		}
		users := r.Form["user"]
		for _, user := range users {
			//get included checkbox value
			includedCheckbox := r.FormValue("includedCheckbox_" + user)
			//Check that the user has been included
			if includedCheckbox != "" {
				log.Println("User: " + user)
				read = "t"
				writeCheckbox := r.FormValue("writeCheckbox_" + user)
				//Check that the user has write privlages
				if writeCheckbox != "" {
					write = "t"
				} else {
					write = "f"
				}

				_, err = stmt.Exec(noteId, user, read, write)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("Read:" + read) //For testing
			}

		}
		log.Println(users) //for testing
		fmt.Fprintf(w, "New Note Added")
		//User is not logged in
	} else {
		//fmt.Fprintf(w, "You are not logged in!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

/*===============================LIST ALL NOTES BELONGING TO USER==================================*/

func listNotes(w http.ResponseWriter, r *http.Request) {
	//Check if user is still online
	if userStillLoggedIn(r) {
		var notes []Note
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			log.Fatal(err)
		}
		username := usernameCookie.Value
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
	//Not 100% sure
	//return listAllUsersSQL(loggedInUser)
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
			http.Redirect(w, r, "/listNotes", http.StatusSeeOther)
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

//==================GET ALL OWNED NOTES=====================================
func getOwndedNotes(username string) (notes []Note) {
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

//==========GET NOTES That you are only appart of====================
//Still need to do
func getPartOfNotes(username string) (notes []Note) {
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

//==============SEARCH FOR A NOTE WITH PARTIAL TEXT SQL==============================
func searchNotePartial(w http.ResponseWriter, r *http.Request) {
	var notes Notes
	//Check if user is still online
	if userStillLoggedIn(r) {

		bodyText := mux.Vars(r)["search"]
		log.Println("Partial String:" + bodyText) //For testing
		//gets notes owned that match the pattern
		notes.OwnedNotes = partialSeachOwnedTitle(bodyText, r)
		//get notes part of that match the pattern
		notes.PartOfNotes = partialSearchPartOfTitle(bodyText, r)
		fmt.Println(notes)
		tpl.ExecuteTemplate(w, "listNotes.gohtml", notes)
		//User is not logged in
	} else {
		//Redirect to splash page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

//===================Delete Specific note=============================

func deleteSpecificNote(r *http.Request) (noteDeleted bool) {

	//get the id of the note the user wants to delete
	//ASK floyd if we ca cut out deleteSpecificNote method and go straight to deleteSpecificNoteSQL
	noteid := mux.Vars(r)["id"]
	username := getUserName(r)

	return deleteSpecificNoteSQL(noteid, username)
}
