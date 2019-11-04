package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	var newUser User
	newUser.UserName = "testUserName"
	newUser.GivenName = "testGivenName"
	newUser.FamilyName = "testFamilyName"
	newUser.Email = "test@test.test"
	newUser.Password = "password"

	var newNote Note
	newNote.NoteTitle = "Test Note Title"
	newNote.NoteBody = "Test Note Body. Partial Search Text"
	newNote.NoteOwner = "testUserName"
	newNote.CreatedDate = "2019-11-03"

	db := connectDatabase()
	if assert.NotNil(t, db) {
		assert.Equal(t, "Added user", addUserSQL(newUser), "Should return added user success message")
		// tests utilizing user start here
		assert.True(t, userNameExists("testUserName"), "Username should exist")
		assert.True(t, validatePass("password", "testUserName"), "User should exist")
		assert.Equal(t, "New Note Added", addNoteSQL(newNote), "Should return success message")
		// tests utilizing note start here
		notes := listAllNotesSQL("testUserName")
		assert.Equal(t, "testUserName", notes[0].NoteOwner, "note should exist")

		// all tests end here
		assert.True(t, deleteAllUserNotesSQL("testUserName"), "should return deleted note. possibly have extra notes in table") // have not tested individual note deletion. need to set up du,mmy tables for that
		assert.True(t, deleteSpecificUserSQL("testUserName"), "Should delete test user")
	}

}
