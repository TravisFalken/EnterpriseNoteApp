package main

import (
	"strconv"
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

	var newNote1 Note
	newNote1.NoteTitle = "Test Note Title"
	newNote1.NoteBody = "Test Note Body. Partial Search Text"
	newNote1.NoteOwner = "testUserName"
	newNote1.CreatedDate = "2019-11-03"

	var newNote2 Note
	newNote2.NoteTitle = "Test Note Title using partial search"
	newNote2.NoteBody = "Test Note Body. Partial Search Text"
	newNote2.NoteOwner = "testUserName"
	newNote2.CreatedDate = "2019-11-03"

	db := connectDatabase()
	defer db.Close()

	// Initial test that database is there and connection can be made
	if assert.NotNil(t, db) {

		// Test Add user, This user is used in testing
		assert.Equal(t, "Added user", addUserSQL(newUser), "Should return added user success message")
		// tests utilizing user start here

		// Test user name can be found on database
		assert.True(t, userNameExists("testUserName"), "Username should exist")

		// test password can be found and validate it for Test user
		assert.True(t, validatePass("password", "testUserName"), "User should exist")

		// Test Add new note. this note will be used in testing
		assert.Equal(t, "New Note Added", addNoteSQL(newNote1), "Should return success message")
		assert.Equal(t, "New Note Added", addNoteSQL(newNote2), "Should return success message")
		// tests utilizing note start here

		// list of notes created for test functions. Testing list all notes
		notes := listAllNotesSQL("testUserName")
		assert.Equal(t, "testUserName", notes[0].NoteOwner, "note should exist with user name")

		// list of all users created for testing. testing that users can be found on data base
		users := listAllUsersSQL("n/a") // na is dummy name added to function
		assert.NotEqual(t, "", users[0], "Should not be null, should have users")

		// note returned from partial search of body. Test that note returned was correct
		partialSearchedNoteBody := partialTextBodySearchSQL("Parti", "testUserName")
		assert.Contains(t, partialSearchedNoteBody[0].NoteBody, "Parti", "note should exist")

		// note return from partial title search. Test returned note was correct
		partialSearchedNoteOwnedTitle := partialSeachOwnedTitleSQL("usi", "testUserName")
		assert.Contains(t, partialSearchedNoteOwnedTitle[0].NoteTitle, "usi", "note should exist")

		// note return for priveleileges of note partial text search. test of part of note search
		partialSearchedNotePartOfTitle := partialSearchPartOfTitleSQL("t", "Trav3")
		assert.Contains(t, partialSearchedNotePartOfTitle[0].NoteTitle, "t", "note should exist")

		// all user and note tests end here
		// Test Delete specific note
		assert.True(t, deleteSpecificNoteSQL(strconv.Itoa(partialSearchedNoteBody[0].NoteID), "testUserName"), "should return deleted note")

		// test delete all remaining user notes
		assert.True(t, deleteAllUserNotesSQL("testUserName"), "should return deleted note")

		// test delete test user
		assert.True(t, deleteSpecificUserSQL("testUserName"), "Should delete test user")
	}

}
