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
	router.HandleFunc("/editNote/{id}", editNote).Methods("GET") // For displaying edit note html page
	//router.HandleFunc("/editNote", updateNote).Methods("POST") //For updating edit note
	router.HandleFunc("/listNotes", allNotes).Methods("GET")
	router.HandleFunc("/addUsers/{id}", listAvaliablePermissions).Methods("GET")
	router.HandleFunc("/addPrivileges/{id}", addPermissions).Methods("POST")
	router.HandleFunc("/listPrivileges/{id}", showPrivileges).Methods("GET")
	router.HandleFunc("/editPrivileges/{id}", editPrivileges).Methods("POST")
	router.HandleFunc("/login", login) //Can be a post and a get method so you know when user is loggin in
	router.HandleFunc("/logout", logout).Methods("GET")
	router.HandleFunc("/deleteNote/{id}", deleteNote).Methods("GET")  //Changed Method from delete to Get because browsers don't support delete method
	router.HandleFunc("/updateNote/{id}", updateNote).Methods("POST") //Updates a note had to change from PUT to POST becasue of formData bug
	router.HandleFunc("/searchNotes", searchNotePartial).Queries("search", "{search}").Methods("GET")
	//router.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	router.PathPrefix("/javascript/").Handler(http.StripPrefix("/javascript/", http.FileServer(http.Dir("./javascript")))) // Handler for serving files within the javascript folder
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))                      //handler for serving files within the css folder

	/////////////////////GROUPS SECTION///////////////////////////
	router.HandleFunc("/viewGroups", viewAllSavedGroups).Methods("GET")
	router.HandleFunc("/createGroup", createGroup).Methods("GET") // For html page
	router.HandleFunc("/addGroup", addGroup).Methods("POST")
	router.HandleFunc("/deleteGroup/{id}", deleteGroup).Methods("GET")
	router.HandleFunc("/viewEditGroupUsers/{id}", viewEditGroupUsers).Methods("GET") // For html page
	router.HandleFunc("/editGroupUsers/{id}", editGroupUsers).Methods("POST")
	router.HandleFunc("/viewGroup/{id}", viewGroup).Methods("GET")
	router.HandleFunc("/updateGroup/{id}", updateGroup).Methods("POST")
	router.HandleFunc("/AddUsersGroup/{id}", viewAddGroupUsers).Methods("GET") // For html page
	router.HandleFunc("/AddUsersGroup/{id}", addGroupUsers).Methods("POST")
	//router.HandleFunc("/editGroup/{id}", editGroup).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
