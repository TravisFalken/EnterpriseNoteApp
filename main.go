package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Note struct {
	NoteID      int       `json: "noteID"`
	NoteTitle   string    `json:"noteTitle"`
	NoteBody    string    `json: "noteBody"`
	CreatedDate time.Time `json: "createdDate"`
	NoteOwner   string    `json:"noteOwner"`
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
	user.UserName = "TravisFalken"
	user.Email = "travis.falkenberg141@gmail.com"
	user.FamilyName = "Falkenberg"
	user.GivenName = "Travis"
	user.Password = "1234"
	note.NoteID = 1
	note.NoteTitle = "test"
	note.NoteBody = "This is a test for the note body"
	note.CreatedDate = time.Now()

	fmt.Println(note)
	setupDB()
	addNote()
	addNote()
	addUser()
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
				createdDate character varying(50) NOT NULL,
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
	getUserName, err := db.Prepare("Select username FROM Client WHERE username = $1")
	if err != nil {
		log.Fatal(err)
	}
	err = getUserName.QueryRow(username).Scan(&name)
	//if error username does not exist in database
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	//Username does exist in the database
	return true
}

//====================ADD USER=====================================
func addUser() {
	fmt.Println("Entered addUser()") // For testing

	var newUser User
	newUser.UserName = "Trav3"
	newUser.Email = "Travis.Falkenberg141@gmail.com"
	newUser.FamilyName = "Falkenberg"
	newUser.GivenName = "Travis"
	newUser.Password = "1234"

	db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if !userNameExists(newUser.UserName) {
		//Prepare insert to stop SQL injections
		log.Println("Entered add user if statement")
		stmt, err := db.Prepare("INSERT INTO Client VALUES($1,$2,$3,$4,$5)")
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(newUser.UserName, newUser.Password, newUser.Email, newUser.GivenName, newUser.FamilyName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Added user")
	} else {
		fmt.Println("Username already exists")
	}

}

//==========================ADD NOTE=============================================================
func addNote() {
	//Make sure user is still logged in
	db, err := sql.Open("postgres", "user=postgres password=password dbname=noteBookApp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//Prepare insert to stop SQL injections
	stmt, err := db.Prepare("INSERT INTO Note (notetitle, notebody, createddate, noteowner) VALUES($1,$2,$3,$4)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(note.NoteTitle, note.NoteBody, note.CreatedDate, user.UserName)
	if err != nil {
		log.Fatal(err)
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
