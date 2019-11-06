package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	assert := assert.New(t)

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
	if assert.NotNil(db) {

		// Test Add user, This user is used in testing
		assert.True(addUserSQL(newUser), "Should return added user success")
		// tests utilizing user start here

		// Test user name can be found on database
		assert.True(userNameExists("testUserName"), "Username should exist")

		// test password can be found and validate it for Test user
		assert.True(validatePass("password", "testUserName"), "User should exist")

		// Test Add new note. this note will be used in testing
		note1ID := addNoteSQL(newNote1)
		assert.NotEmpty(note1ID, "Should return note id")
		note2ID := addNoteSQL(newNote2)
		assert.NotEmpty(note2ID, "Should return note id")
		// tests utilizing note start here

		// list of notes created for test functions. Testing list all notes
		notes := listAllNotesSQL("testUserName")
		assert.Equal("testUserName", notes[0].NoteOwner, "note should exist with user name")

		// list of all users created for testing. testing that users can be found on data base
		users := listAllUsersSQL("n/a") // na is dummy name added to function
		assert.NotEqual("", users[0], "Should not be null, should have users")

		// note returned from partial search of body. Test that note returned was correct
		partialSearchedNoteBody := partialTextBodySearchSQL("Parti", "testUserName")
		assert.Contains(partialSearchedNoteBody[0].NoteBody, "Parti", "note should exist")

		// note return from partial title search. Test returned note was correct
		partialSearchedNoteOwnedTitle := partialSeachOwnedTitleSQL("usi", "testUserName")
		assert.Contains(partialSearchedNoteOwnedTitle[0].NoteTitle, "usi", "note should exist")

		// note return for priveleileges of note partial text search. test of part of note search
		partialSearchedNotePartOfTitle := partialSearchPartOfTitleSQL("t", "Trav3")
		assert.Contains(partialSearchedNotePartOfTitle[0].NoteTitle, "t", "note should exist")

		// notes returned for owned notes. Test owned notes have correct owner
		ownedNotes := getOwnedNotesSQL("testUserName")
		assert.Equal("testUserName", ownedNotes[0].NoteOwner, "Note should have owner of testUserName")

		// all user and note tests end here
		// Test Delete specific note
		assert.True(deleteSpecificNoteSQL(strconv.Itoa(partialSearchedNoteBody[0].NoteID), "testUserName"), "should return deleted note")

		// test delete all remaining user notes
		assert.True(deleteAllUserNotesSQL("testUserName"), "should return deleted note")

		// test delete test user
		assert.True(deleteSpecificUserSQL("testUserName"), "Should delete test user")

	}

}
