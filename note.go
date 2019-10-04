package main

type Note struct {
	NoteID      int    `json: "noteID"`
	NoteTitle   string `json:"noteTitle"`
	NoteBody    string `json: "noteBody"`
	CreatedDate string `json: "createdDate"`
	NoteOwner   string `json:"noteOwner"`
}
