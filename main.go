package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

var note Note //For testing dummy data
var user User //For testing dummy data

func main() {

	fmt.Println(note)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/home", home).Methods("GET")
	router.HandleFunc("/createUser", addUser).Methods("POST")
	router.HandleFunc("/createNote", createNote).Methods("GET") //For the html page
	router.HandleFunc("/addNote", addNote).Methods("POST")
	router.HandleFunc("/signUp", signUp).Methods("GET")
	router.HandleFunc("/listAllNotes", listNotes).Methods("GET")
	router.HandleFunc("/listNotes", allNotes).Methods("GET")
	router.HandleFunc("/login", login) //Can be a post and a get method so you know when user is loggin in
	router.HandleFunc("/logout", logout).Methods("GET")
	router.HandleFunc("/deleteNote/{id}", deleteNote).Methods("DELETE")
	router.HandleFunc("/searchNotes/{id}", searchNotePartial).Methods("GET")
	//router.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	log.Fatal(http.ListenAndServe(":8080", router))
}
