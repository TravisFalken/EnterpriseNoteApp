package model

type Note struct {
	NoteID int    `json:"noteID"`
	UserID int    `json:"userID"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Date   string `json:"Date"` // may change to datetime or just parse it
}
