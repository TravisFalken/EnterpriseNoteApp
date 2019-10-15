package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// type Note struct {
// 	NoteID      int    `json: "noteID"` // can remove but left here just incase we need to remeber anything
// 	NoteTitle   string `json:"noteTitle"`
// 	NoteBody    string `json: "noteBody"`
// 	CreatedDate string `json: "createdDate"`
// 	NoteOwner   string `json:"noteOwner"`
// }

// type User struct {
// 	UserName   string `json:"userName"`
// 	Password   string `json:"password"` //This will normally be encripted
// 	Email      string `json:"email"`
// 	GivenName  string `json:"givenName"`
// 	FamilyName string `json:"familyName"`
// 	SessionID  string `json:"sessionID"`
//}

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
	router.HandleFunc("/createUser", addUser).Methods("POST")
	router.HandleFunc("/createNote", addNote).Methods("POST")
	router.HandleFunc("/listAllNotes", listNotes).Methods("GET")
	router.HandleFunc("/login", login).Methods("GET")
	router.HandleFunc("/logout", logout).Methods("GET")
	router.HandleFunc("/deleteNote/{id}", deleteNote).Methods("DELETE")
	router.HandleFunc("/searchNotes/{id}", searchNotePartial).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
