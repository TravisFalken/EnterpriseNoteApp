package main

type Note struct {
	NoteID      int    `json: "noteID"`
	NoteTitle   string `json:"noteTitle"`
	NoteBody    string `json: "noteBody"`
	CreatedDate string `json: "createdDate"`
	NoteOwner   string `json:"noteOwner"`
	Read        string `json:"read"`  //for when user is part of note
	Write       string `json:"write"` //for when user is part of note
}
