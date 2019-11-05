package main

import (
	"database/sql"
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

//Edit privileges page for a note
func showPrivileges(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		noteid := getNoteID(r)
		privileges := struct {
			NoteID     string
			Privileges []privlige
		}{
			NoteID:     noteid,
			Privileges: getNotePrivileges(noteid),
		}

		tpl.ExecuteTemplate(w, "editPrivileges.gohtml", privileges)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

//Add people to existing note
func listAvaliablePermissions(w http.ResponseWriter, r *http.Request) {

	//Make sure user is still logged in
	if userStillLoggedIn(r) {
		//validate that user is note owner
		if noteOwner(r) {
			//create a temp struct to hold note id and usernames
			tempStruct := struct {
				NoteID string
				Users  []string
			}{
				NoteID: getNoteID(r),
				Users:  getAvaliableUsers(r),
			}

			tpl.ExecuteTemplate(w, "addPermissions.gohtml", tempStruct)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
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
				log.Println("User: " + user) //for testing
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
//===========NOT USING CAN DELETE==========
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
		notes.OwnedNotes = getOwnedNotesSQL(username)
		notes.PartOfNotes = getPartOfNotesSQL(username)
		log.Println(notes) //For testing
		tpl.ExecuteTemplate(w, "listNotes.gohtml", notes)

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//==============SEARCH FOR A NOTE WITH PARTIAL TEXT SQL==============================
func searchNotePartial(w http.ResponseWriter, r *http.Request) {
	var notes Notes
	//Check if user is still online
	if userStillLoggedIn(r) {

		bodyText := mux.Vars(r)["search"]
		username := getUserName(r)
		log.Println("Partial String:" + bodyText) //For testing
		//gets notes owned that match the pattern
		notes.OwnedNotes = partialSeachOwnedTitleSQL(bodyText, username)
		//get notes part of that match the pattern
		notes.PartOfNotes = partialSearchPartOfTitleSQL(bodyText, username)
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
	noteid := getNoteID(r)
	username := getUserName(r)
	log.Println("NoteID + Username : " + noteid + username) //for testing

	return deleteSpecificNoteSQL(noteid, username)
}

//===================Update a Note=======================

func updateNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered update note") // for testing
	if userStillLoggedIn(r) {
		//validate if note owner
		if noteOwner(r) {
			//update note
			if updateOwnedNote(r) {
				http.Redirect(w, r, "/listNotes", http.StatusSeeOther)
			} else {
				http.Error(w, "Could not update Note!", http.StatusExpectationFailed)
			}
		}
		//validate if has write access
		if checkWritePermissions(r) {
			if updatePartOfNote(r) {
				http.Redirect(w, r, "/listNotes", http.StatusSeeOther)
			} else {
				http.Error(w, "Could not update Note!", http.StatusExpectationFailed)
			}
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//check if user has read permissions
func readPermissions(r *http.Request) (readPremission bool) {
	username := getUserName(r)
	var read string
	//Get note id from http
	noteid := getNoteID(r)

	return readPermissionsSQL(username, noteid, read)

}

//Check if user has write permissions
func checkWritePermissions(r *http.Request) (writePermission bool) {
	username := getUserName(r)
	var write string

	noteid := getNoteID(r)

	return checkWritePermissionsSQL(username, noteid, write)

}

//Check if user is a owner of a note
func noteOwner(r *http.Request) bool {
	var owner string
	username := getUserName(r)
	noteid := getNoteID(r)

	return noteOwnerSQL(username, noteid, owner)

}

//get all users that are not already part of note
func getAvaliableUsers(r *http.Request) (users []string) {

	username := getUserName(r)
	noteid := getNoteID(r)

	return getAvaliableUsersSQL(username, noteid)
}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR SQL TO DATABASE
// ------------------------------------------------------------------------------------------------------------------------------------------------

// TRAVIS. i dont understand what is happening here. and it seems either way you will return the same note. weather they pass the if statement or not

//get note based on note id and user permissions
func getPartOfNote(r *http.Request) (note Note) {
	//Check to see if user has read permission or is note owner
	if readPermissions(r) {
		//Connect to database
		db := connectDatabase()
		defer db.Close()
		//Get username
		username := getUserName(r)
		noteid := getNoteID(r)
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
// CLEAR SQL TO DATABASE
// ------------------------------------------------------------------------------------------------------------------------------------------------

func getOwnedNote(r *http.Request) (note Note) {
	if noteOwner(r) {
		//Connect to database
		db := connectDatabase()
		defer db.Close()
		//get username
		username := getUserName(r)
		noteid := getNoteID(r)
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
// CLEAR SQL TO DATABASE
// ------------------------------------------------------------------------------------------------------------------------------------------------

//Update note that you own
func updateOwnedNote(r *http.Request) (success bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//get username
	username := getUserName(r)
	//get the note id that we need to update
	noteid := getNoteID(r)
	//get values from form
	title := r.FormValue("title")
	body := r.FormValue("body")
	fmt.Println("This is the title and body of update: " + title + " " + body) //for testing
	if updateOwnedNoteSQL(title, body, noteid, username) {
		return true
	} else {
		return false
	}

}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR SQL TO DATABASE
// ------------------------------------------------------------------------------------------------------------------------------------------------

//Update note you are part of
func updatePartOfNote(r *http.Request) bool {
	if userStillLoggedIn(r) {

		//get username
		//username := getUserName(r)
		//get note id
		noteID := getNoteID(r)

		//Get value from form
		body := r.FormValue("body")
		fmt.Println("This is the update body: " + body) //This is for testing

		if updatePartOfNoteSQL(noteID, body) {
			return true
		}
	}
	return false
}

//Need to move to database.go but just putting it here for use sake
//Add more users to a note
func addPermissions(w http.ResponseWriter, r *http.Request) {
	log.Println("Entered add premission") //For testing
	if userStillLoggedIn(r) {

		noteid := getNoteID(r)

		var read string
		var write string
		//DO NOT REMOVE PRINTLN !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		log.Println("value: " + r.FormValue("includedCheckbox_Vaughn1")) //For tessting

		users := r.Form["user"]
		for _, user := range users {
			//get included checkbox value
			log.Println("User:" + user) //For testing
			includedCheckbox := r.FormValue("includedCheckbox_" + user)
			log.Println("Included Checkbox: " + includedCheckbox) // for testing
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
				//Add to database
				if !addPermissionSQL(noteid, user, read, write) {
					http.Error(w, "Database Error", http.StatusInternalServerError)
				}

				log.Println("Read:" + read) //For testing
			}

		}

		http.Redirect(w, r, "/listNotes", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func editPrivileges(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) && noteOwner(r) {
		var included string
		var writeValue string
		var write string
		//DO NOT REMOVE PRINTLN !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		log.Println("value: " + r.FormValue("includedCheckbox_Vaughn1")) //For tessting
		noteid := getNoteID(r)
		users := r.Form["user"]
		log.Println(users[0])
		for _, user := range users {
			included = r.FormValue("includedCheckbox_" + user)
			//If user not included anymore note privilege will be moved
			if included == "" {
				removePrivilege(noteid, user)
				continue
			}
			writeValue = r.FormValue("writeCheckbox_" + user)
			if writeValue != "" {
				write = "t"
			} else {
				write = "f"
			}
			updatePrivilege(noteid, user, write)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//get note id from url
func getNoteID(r *http.Request) (noteid string) {
	noteid = mux.Vars(r)["id"]
	return noteid
}
