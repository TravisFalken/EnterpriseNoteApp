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
	setupDB()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/createUser", addUser).Methods("POST")
	router.HandleFunc("/createNote", addNote).Methods("POST")
	router.HandleFunc("/listAllNotes", listNotes).Methods("GET")
	router.HandleFunc("/login", login).Methods("GET")
	router.HandleFunc("/logout", logout).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

//Set up database
func setupDB() {
	db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}
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
	/*
		userTableQuery := `DROP TABLE IF EXISTS Client;
				CREATE TABLE Client(
					userName  character varying(50) NOT NULL PRIMARY KEY,
					password character varying(50) NOT NULL,
					email character varying(250) NOT NULL,
					givenName character varying(50) NOT NULL,
					familyName character varying(50) NOT NULL,
					sessionID character varying(250)
				);`
	*/
	noteTableQuery := `DROP TABLE IF EXISTS Note;
				CREATE TABLE Note(
				noteID serial PRIMARY KEY,
				noteTitle character varying(50) NOT NULL,
				noteBody TEXT NOT NULL,
				createdDate character varying(250) NOT NULL,
				noteOwner character varying(50) NOT NULL,
				FOREIGN KEY(noteOwner) REFERENCES Client(userName)
			);`
	/*
		_, err = db.Exec(userTableQuery)
		if err != nil {
			log.Fatal(err)
		}
	*/
	_, err = db.Exec(noteTableQuery)
	if err != nil {
		log.Fatal(err)
	}

}

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
	db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Prepare query to check if the username already exists
	getUserName, err := db.Prepare("Select username FROM _user WHERE username = $1")
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
	db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//prepare statement to check for password
	stmt, err := db.Prepare("SELECT password FROM _user WHERE password = $1")
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

	//Create connection to server
	db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if !userNameExists(newUser.UserName) {
		//Prepare insert to stop SQL injections
		log.Println("Entered add user if statement")
		stmt, err := db.Prepare("INSERT INTO _user VALUES($1,$2,$3,$4,$5)")
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(newUser.UserName, newUser.Password, newUser.Email, newUser.GivenName, newUser.FamilyName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Added user")
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
		//Get the body and put it into a a note struct
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &newNote)

		//set the created date of the note
		newNote.CreatedDate = noteTime.String()
		newNote.NoteOwner = user.UserName

		db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		//Prepare insert to stop SQL injections
		stmt, err := db.Prepare("INSERT INTO _note (notetitle, notebody, createddate, noteowner) VALUES($1,$2,$3,$4)")
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(newNote.NoteTitle, newNote.NoteBody, newNote.CreatedDate, user.UserName)
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

	var notes []Note
	//Connect to DB
	db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	stmt, err := db.Prepare("SELECT * FROM _note WHERE noteowner=$1")
	var note Note
	rows, err := stmt.Query(user.UserName)
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
}

//============================================LIST ALL REGISTERED USERS=====================
//This function creates a list of all registerd users
func listAllUses(loggedInUser string) []string {
	var users []string
	var username string
	//Connect to DB
	db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	//Send query to the db
	rows, err := db.Query("SELECT username FROM _user")
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
		fmt.Fprintf(w, "Successfully logged out")
	} else {
		fmt.Fprintf(w, "Already Logged out")
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
