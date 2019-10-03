package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Note struct {
	NoteID      int    `json: "noteID"`
	NoteTitle   string `json:"noteTitle"`
	NoteBody    string `json: "noteBody"`
	CreatedDate string `json: "createdDate"`
	NoteOwner   string `json:"noteOwner"`
}

type User struct {
	UserName   string `json:"userName"`
	Password   string `json:"password"` //This will normally be encripted
	Email      string `json:"email"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	SessionID  string `json:"sessionID"`
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

var note Note //For testing dummy data
var user User //For testing dummy data

func main() {

	//Ask john if I should make db global so I do not have to connect to it everytime I query
	user.UserName = "TravisFalken"
	user.Email = "travis.falkenberg141@gmail.com"
	user.FamilyName = "Falkenberg"
	user.GivenName = "Travis"
	user.Password = "1234"

	fmt.Println(note)
	//setupDB()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/home", home).Methods("GET")
	router.HandleFunc("/createUser", addUser).Methods("POST")
	router.HandleFunc("/createNote", addNote).Methods("POST")
	router.HandleFunc("/signUp", signUp).Methods("GET")
	router.HandleFunc("/listAllNotes", listNotes).Methods("GET")
	router.HandleFunc("/login", login) //Can be a post and a get method so you know when user is loggin in
	router.HandleFunc("/logout", logout).Methods("GET")
	router.HandleFunc("/deleteNote/{id}", deleteNote).Methods("DELETE")
	router.HandleFunc("/searchNotes/{id}", searchNotePartial).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

//Set up database
func setupDB() {
	//Connect to db
	db := connectDatabase()
	defer db.Close()

	/*
		test := `DROP TABLE IF EXISTS Client;
				CREATE TABLE Client(
				Year integer,
				Measure character varying(50),
				Technology character varying(50),
				Value integer,
				ValueUnit character varying(50),
				ValueLabel character varying(50),
				NullReason character varying(50)
			);`
	*/

	_note_privilegesQuery := `CREATE TABLE _note_privileges( -- added underscore here to keep naming convention
				note_privileges_id integer PRIMARY KEY NOT NULL,
				note_id integer,
				user_name character varying(50),
				read CHAR(1), -- t for true  f for false
				write CHAR(1) -- t for true  f for false
			);`

	userTableQuery := `DROP TABLE IF EXISTS _user CASCADE;
				CREATE TABLE _user(
					user_name  character varying(50) NOT NULL PRIMARY KEY,
					password character varying(50) NOT NULL,
					email character varying(250) NOT NULL,
					given_name character varying(50) NOT NULL,
					family_name character varying(50) NOT NULL,
					session_id character varying(250)
				);`
	noteTableQuery := `DROP TABLE IF EXISTS _Note;
				CREATE TABLE _Note(
				note_id serial PRIMARY KEY,
				title character varying(50) NOT NULL,
				body TEXT NOT NULL,
				date_created character varying(250) NOT NULL,
				note_owner character varying(50) NOT NULL,
				FOREIGN KEY(note_owner) REFERENCES _user(user_name)
			);`

	/*
		_, err = db.Exec(userTableQuery)
		if err != nil {
			log.Fatal(err)
		}
	*/
	_, err := db.Exec(userTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(noteTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(_note_privilegesQuery)
}

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
		err := tpl.ExecuteTemplate(w, "home.gohtml", nil)
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
				//Send the home page to user
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
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &newNote)

		//set the created date of the note
		newNote.CreatedDate = noteTime.String()
		newNote.NoteOwner = username

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

		fmt.Fprintf(w, "New Note Added")
		//User is not logged in
	} else {
		fmt.Fprintf(w, "You are not logged in!")
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
		var notes []Note
		username := getUserName(r)
		//Connect to db
		db := connectDatabase()
		defer db.Close()
		bodyText := mux.Vars(r)["id"]
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
		fmt.Println(notes)
		//just sending it straight to frontend for testing
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
		//User is not logged in
	} else {
		fmt.Fprintf(w, "Not Logged in!")
	}

}

//insertion sort
/*
func insertionSort(arr []Student) []Student {
	for i := 1; i < len(arr); i++ {
		key := arr[i].LastName
		ts := arr[i]
		j := i - 1
		for j >= 0 && key < arr[j].LastName {
			arr[j+1] = arr[j]
			j -= 1
		}
		arr[j+1] = ts
	}

	return arr
}

//Binary Search
func binarySearch(arr []Student, inputName string, inputSurname string) (result Student) {
	low := 0
	high := len(arr) - 1
	mid := 0
	var mid_value Student

	for low <= high {
		mid = low + (high-low)/2 //middle of the array
		mid_value = arr[mid]

		if mid_value.LastName == inputSurname {
			if mid_value.FirstName == inputName {
				return arr[mid] //return the found result
			}
		} else if mid_value.LastName <
			inputSurname {
			low = mid + 1 //left/lower side of the middle
		} else {
			high = mid - 1 //right/upper side of the middle
		}
	}

	return result //Not found so return no position
}
*/
