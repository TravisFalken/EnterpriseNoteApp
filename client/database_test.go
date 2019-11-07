package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	assert := assert.New(t)

	var newUser1 User
	newUser1.UserName = "testUserName1"
	newUser1.GivenName = "testGivenName"
	newUser1.FamilyName = "testFamilyName"
	newUser1.Email = "test@test.test"
	newUser1.Password = "password"

	var newUser2 User
	newUser2.UserName = "testUserName2"
	newUser2.GivenName = "testGivenName2"
	newUser2.FamilyName = "testFamilyName2"
	newUser2.Email = "test@test.test"
	newUser2.Password = "password"

	var newNote1 Note
	newNote1.NoteTitle = "Test Note Title"
	newNote1.NoteBody = "Test Note Body. Partial Search Text"
	newNote1.NoteOwner = "testUserName1"
	newNote1.CreatedDate = "2019-11-03"

	var newNote2 Note
	newNote2.NoteTitle = "Test Note Title using partial search"
	newNote2.NoteBody = "Test Note Body. Partial Search Text"
	newNote2.NoteOwner = "testUserName1"
	newNote2.CreatedDate = "2019-11-03"

	var newGroup Group
	newGroup.GroupTitle = "testGroup"
	newGroup.GroupOwner = newUser1.UserName
	newGroup.GroupRead = "t"
	newGroup.GroupWrite = "t"

	db := connectDatabase()
	defer db.Close()

	// -------------------------------------------------------------
	// Database testing. 100% function coverage
	// -------------------------------------------------------------

	// Initial test that database is there and connection can be made
	if assert.NotNil(db) {

		// Test Add user, This user is used in testing
		assert.True(addUserSQL(newUser1), "Should return added user success")
		assert.True(addUserSQL(newUser2), "Should return added user success")
		// tests utilizing user start here -----------------------------------------------

		// Test user name can be found on database
		assert.True(userNameExists("testUserName1"), "Username should exist")

		// test password can be found and validate it for Test user
		assert.True(validatePass("password", "testUserName1"), "User should exist")

		// Test Add new note. this note will be used in testing
		note1ID := addNoteSQL(newNote1)
		assert.NotEmpty(note1ID, "Should return note id")
		note2ID := addNoteSQL(newNote2)
		assert.NotEmpty(note2ID, "Should return note id")
		// tests utilizing note start here --------------------------------------------------

		// Test adding user2 to note 1
		assert.True(addPermissionSQL(note1ID, newUser2.UserName, "t", "t"), "Permissions should be added on first note")

		// Test adding user2 to note 2
		assert.True(addPermissionSQL(note2ID, newUser2.UserName, "t", "f"), "Permissions should be added on second note")
		// tests using permissions start here -----------------------------------------------

		// list of notes created for test functions. Testing list all notes
		notes := listAllNotesSQL("testUserName1")
		assert.Equal("testUserName1", notes[0].NoteOwner, "note should exist with user name")

		// list of all users created for testing. testing that users can be found on data base
		users := listAllUsersSQL("n/a") // na is dummy name added to function
		assert.NotEqual("", users[0], "Should not be null, should have users")

		// note returned from partial search of body. Test that note returned was correct
		partialSearchedNoteBody := partialTextBodySearchSQL("Parti", "testUserName1")
		assert.Contains(partialSearchedNoteBody[0].NoteBody, "Parti", "note should exist")

		// note return from partial title search. Test returned note was correct
		partialSearchedNoteOwnedTitle := partialSeachOwnedTitleSQL("usi", "testUserName1")
		assert.Contains(partialSearchedNoteOwnedTitle[0].NoteTitle, "usi", "note should exist")

		// note return for priveleileges of note partial text search. test of part of note search
		partialSearchedNotePartOfTitle := partialSearchPartOfTitleSQL("t", "Trav3")
		assert.Contains(partialSearchedNotePartOfTitle[0].NoteTitle, "t", "note should exist")

		// notes returned for owned notes. Test owned notes have correct owner
		ownedNotes := getOwnedNotesSQL("testUserName1")
		assert.Equal("testUserName1", ownedNotes[0].NoteOwner, "Note should have owner of testUserName1")

		// Test part of note SQL
		partOfNote := getPartOfNoteSQL(note1ID, newUser2.UserName)
		assert.Equal(partOfNote.NoteOwner, newUser1.UserName, "Note returned should be from user 1")

		// Test read permissions SQL
		assert.True(readPermissionsSQL(newUser2.UserName, note1ID, ""), "should return true for read permissions")

		// Test write permissions SQL
		assert.False(checkWritePermissionsSQL(newUser2.UserName, note2ID, ""), "Should return false for write permissions")

		// Test note owner sql.
		assert.True(noteOwnerSQL(newUser1.UserName, note1ID, ""), "SHould return true for note owner")

		// Test Available user
		availableUsers := getAvaliableUsersSQL(newUser1.UserName, note1ID)
		assert.NotEqual(availableUsers[0], newUser1.UserName, "returned users should not contain same user")

		// Test get not priviliges
		priviligesRecieved := getNotePrivileges(note1ID)
		assert.Equal(priviligesRecieved[0].Username, newUser2.UserName, "User with priviliges should be new user 2")

		// Test update priviliges
		assert.True(updatePrivilege(note2ID, newUser2.UserName, "t"), "Write should be updated to true")

		// Test Update Owned note
		assert.True(updateOwnedNoteSQL("Test update", "Test Update", note1ID, newUser1.UserName), "Should be able to update as user1 is owner of this note")

		// Test update part of note
		assert.True(updatePartOfNoteSQL(note1ID, "Test Update"), "Should be able to update")

		// Test Create Group
		groupID := createNewGroup(newGroup.GroupTitle, newGroup.GroupOwner, newGroup.GroupRead, newGroup.GroupWrite)
		assert.NotEmpty(groupID, "should hold created group id")
		// Group tests start here --------------------------------------------------------------------------------------------

		//

		// all users, groups, notes and permissions tests end here ------------------------------------------------------------

		// Test remove group
		assert.True(removeGroup(groupID), "Should remove group")

		// Test remove permissions
		assert.True(removePrivilege(note1ID, newUser2.UserName), "permissions should be removed on note 1")
		assert.True(removePrivilege(note2ID, newUser2.UserName), "permissions should be removed on note 2")

		// Test Delete specific note
		assert.True(deleteSpecificNoteSQL(strconv.Itoa(partialSearchedNoteBody[0].NoteID), "testUserName1"), "should return deleted note")

		// test delete all remaining user notes
		assert.True(deleteAllUserNotesSQL("testUserName1"), "should return deleted note")

		// test delete test user
		assert.True(deleteSpecificUserSQL("testUserName1"), "Should delete test user")
		assert.True(deleteSpecificUserSQL("testUserName2"), "Should delete test user")

	}

}
