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
	db := connectDatabase()
	if assert.NotNil(t, db) {
		assert.Equal(t, "Added user", addUserSQL(newUser), "Should return added user success message")
		assert.True(t, userNameExists("testUserName"), "Username should exist")
		assert.True(t, validatePass("password", "testUserName"), "User should exist")

		assert.True(t, deleteSpecificUserSQL("testUserName"), "Should delete test user")
	}

}
