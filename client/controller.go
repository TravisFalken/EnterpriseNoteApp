package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//This method runs the first time the user trys to access the webapp
//Looking for a better function name XD
func index(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//============This is the home page===============================
func home(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		sessionid := getSession(r)
		user := getUser(sessionid)
		err := tpl.ExecuteTemplate(w, "home.gohtml", user)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//===========================THE SIGNUP Page============================================
func signUp(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	} else {
		tpl.ExecuteTemplate(w, "signup.gohtml", nil)
	}

}

//=====================THE CREATE NOTE PAGE===========================================
func createNote(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		user := getUserName(r)
		users := listAllUsersSQL(user)
		groups := getAllGroups(user)
		//Create temp struct
		groupsAndUsers := struct {
			Users  []string
			Groups []Group
		}{
			Users:  users,
			Groups: groups,
		}

		tpl.ExecuteTemplate(w, "createNote.gohtml", groupsAndUsers)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//====================EDIT NOTE PAGE================================================
func editNote(w http.ResponseWriter, r *http.Request) {
	//Make sure user is still logged in
	if userStillLoggedIn(r) {
		//validate that user can actually read the note
		if readPermissions(r) || noteOwner(r) {
			var note Note
			//get note if user is part of
			if readPermissions(r) {
				note = getPartOfNote(r)
				//get note if user us owner
			} else {
				note = getOwnedNote(r)
			}
			//get the note
			tpl.ExecuteTemplate(w, "editNote.gohtml", note)
		} else {
			http.Redirect(w, r, "/listNotes", http.StatusSeeOther)
		}

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Edit privileges page for a note
func showPrivileges(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		noteid := getID(r)
		privileges := struct {
			NoteID     string
			Privileges []Privlige
		}{
			NoteID:     noteid,
			Privileges: getNotePrivileges(noteid),
		}

		tpl.ExecuteTemplate(w, "editPrivileges.gohtml", privileges)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

//Add people to existing note
func listAvaliablePermissions(w http.ResponseWriter, r *http.Request) {

	//Make sure user is still logged in
	if userStillLoggedIn(r) {
		//validate that user is note owner
		if noteOwner(r) {
			//create a temp struct to hold note id and usernames
			tempStruct := struct {
				NoteID string
				Users  []string
			}{
				NoteID: getID(r),
				Users:  getAvaliableUsers(r),
			}

			tpl.ExecuteTemplate(w, "addPermissions.gohtml", tempStruct)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//=========================Checks if user login details are correct=========================================
func login(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {

		var loginUser User
		loginUser.UserName = r.FormValue("username")
		loginUser.Password = r.FormValue("password")
		//Validate that user has typed in input
		if validateInput(loginUser.UserName) && validateInput(loginUser.Password) {
			//make sure username doesn't exist
			if userNameExists(loginUser.UserName) {
				if validatePass(loginUser.Password, loginUser.UserName) {

					//create new session id
					sessionid := newSessionid()
					//Add session id to the user in the database
					addSessionToUser(loginUser, sessionid)
					//set the clients session cookie
					sessionCookie := &http.Cookie{
						Name:  "session",
						Value: sessionid,
					}
					http.SetCookie(w, sessionCookie)
					//Create username cookie
					usernameCookie := &http.Cookie{
						Name:  "username",
						Value: loginUser.UserName,
					}
					//Set a cookie to username
					http.SetCookie(w, usernameCookie)
					//Send the home page to user
					http.Redirect(w, r, "/", http.StatusSeeOther)
				} else {
					http.Error(w, "username and/or password does not match", http.StatusForbidden)
					return
				}
			} else {
				http.Error(w, "username and/or password does not match", http.StatusForbidden)
				return
			}
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//====================ADD USER=====================================
func addUser(w http.ResponseWriter, r *http.Request) {

	var newUser User

	newUser.UserName = r.FormValue("user_name")
	newUser.GivenName = r.FormValue("given_name")
	newUser.FamilyName = r.FormValue("family_name")
	newUser.Email = r.FormValue("email")
	newUser.Password = r.FormValue("password")
	//Validate that the user has filled in the required fields
	if validateInput(newUser.UserName) && validateInput(newUser.Email) && validateInput(newUser.Password) {
		//Validate that username does not already exist
		if !userNameExists(newUser.UserName) {
			//Add user to database
			if addUserSQL(newUser) {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				fmt.Fprintf(w, "Failed to create new user!")
			}
		} else {
			fmt.Fprintf(w, "Username already exists!")
		}
	} else {
		http.Redirect(w, r, "/signUp", http.StatusSeeOther)
	}

}

//==========================ADD NOTE=============================================================
func addNote(w http.ResponseWriter, r *http.Request) {
	//Check if user is logged in
	if userStillLoggedIn(r) {
		var newNote Note
		var noteTime = time.Now()
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			log.Fatal(err)
		}
		username := usernameCookie.Value

		//Get the body and put it into a a note struct
		newNote.NoteTitle = r.FormValue("title")
		newNote.NoteBody = r.FormValue("body")
		newNote.CreatedDate = noteTime.Format("2006-01-02")
		newNote.NoteOwner = username
		//Validate that user has typed in required inputs
		if validateInput(newNote.NoteTitle) && validateInput(newNote.NoteBody) {
			//Add note to database and get noteid back
			noteid := addNoteSQL(newNote)
			//Get users attached to note and add them to the database
			var read string
			var write string
			//validate if user wants to user groups or manually enter users
			useGroup := r.FormValue("useSavedGroup")
			if useGroup == "" {
				users := r.Form["user"]
				for _, user := range users {
					//get included checkbox value
					includedCheckbox := r.FormValue("includedCheckbox_" + user)
					//Check that the user has been included
					if includedCheckbox != "" {
						read = "t"
						writeCheckbox := r.FormValue("writeCheckbox_" + user)
						//Check that the user has write privlages
						if writeCheckbox != "" {
							write = "t"
						} else {
							write = "f"
						}

						//Add permission to the database
						addPermissionSQL(noteid, user, read, write)
					}

				}
			} else {
				groupid := r.FormValue("group")
				//Validate that user is group owner
				if validateGroupOwner(username, groupid) {
					//Gett all of users from group
					users := getGroupUsers(groupid)
					//get write and read premissions from the group
					read, write = getGroupPrivileges(groupid)
					for _, user := range users {
						addPermissionSQL(noteid, user, read, write)
					}
				} else {
					http.Redirect(w, r, "/", http.StatusSeeOther)
				}
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			// fmt.Fprintf(w, "New Note Added")
			//User is not logged in
		} else {
			http.Redirect(w, r, "/createNote", http.StatusSeeOther)
		}
	} else {
		//fmt.Fprintf(w, "You are not logged in!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

//=========================USER LOGOUT======================================================
func logout(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		deleteSesion(w, r)
		//fmt.Fprintf(w, "Successfully logged out")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		//fmt.Fprintf(w, "Already logged out")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

//==============DELETE A NOTE==============================
func deleteNote(w http.ResponseWriter, r *http.Request) {
	//Check if user is still logged in
	if userStillLoggedIn(r) {
		//deletes a note and returns true if deleted
		if deleteSpecificNote(r) {
			http.Redirect(w, r, "/listNotes", http.StatusSeeOther)
		} else {
			fmt.Fprintf(w, "Not Successful")
		}
	} else {
		fmt.Fprintf(w, "You cannot delete note because you are not logged in")
	}
}

//==============List All Notes that user owners and is part of==================================
func allNotes(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		var username = getUserName(r)
		var notes Notes
		notes.OwnedNotes = getOwnedNotesSQL(username)
		notes.PartOfNotes = getPartOfNotesSQL(username)
		tpl.ExecuteTemplate(w, "listNotes.gohtml", notes)

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//==============SEARCH FOR A NOTE WITH PARTIAL TEXT SQL==============================
func searchNotePartial(w http.ResponseWriter, r *http.Request) {
	var notes Notes
	//Check if user is still online
	if userStillLoggedIn(r) {

		bodyText := mux.Vars(r)["search"]
		username := getUserName(r)
		//gets notes owned that match the pattern
		notes.OwnedNotes = partialSeachOwnedTitleSQL(bodyText, username)
		//get notes part of that match the pattern
		notes.PartOfNotes = partialSearchPartOfTitleSQL(bodyText, username)
		tpl.ExecuteTemplate(w, "listNotes.gohtml", notes)
		//User is not logged in
	} else {
		//Redirect to splash page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

//===================Delete Specific note=============================

func deleteSpecificNote(r *http.Request) (noteDeleted bool) {

	//get the id of the note the user wants to delete
	noteid := getID(r)
	username := getUserName(r)

	return deleteSpecificNoteSQL(noteid, username)
}

//===================Update a Note=======================

func updateNote(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		//validate if note owner
		if noteOwner(r) {
			//update note
			if updateOwnedNote(r) {
				http.Redirect(w, r, "/listNotes", http.StatusSeeOther)
			} else {
				http.Error(w, "Could not update Note!", http.StatusExpectationFailed)
			}
		}
		//validate if has write access
		if checkWritePermissions(r) {
			if updatePartOfNote(r) {
				http.Redirect(w, r, "/listNotes", http.StatusSeeOther)
			} else {
				http.Error(w, "Could not update Note!", http.StatusExpectationFailed)
			}
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//check if user has read permissions
func readPermissions(r *http.Request) (readPremission bool) {
	username := getUserName(r)
	var read string
	//Get note id from http
	noteid := getID(r)

	return readPermissionsSQL(username, noteid, read)

}

//Check if user has write permissions
func checkWritePermissions(r *http.Request) (writePermission bool) {
	username := getUserName(r)
	var write string

	noteid := getID(r)

	return checkWritePermissionsSQL(username, noteid, write)

}

//Check if user is a owner of a note
func noteOwner(r *http.Request) bool {
	var owner string
	username := getUserName(r)
	noteid := getID(r)

	return noteOwnerSQL(username, noteid, owner)

}

//get all users that are not already part of note
func getAvaliableUsers(r *http.Request) (users []string) {

	username := getUserName(r)
	noteid := getID(r)

	return getAvaliableUsersSQL(username, noteid)
}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR SQL TO DATABASE
// ------------------------------------------------------------------------------------------------------------------------------------------------

// TRAVIS. i dont understand what is happening here. and it seems either way you will return the same note. weather they pass the if statement or not

//get note based on note id and user permissions
func getPartOfNote(r *http.Request) (note Note) {
	//Check to see if user has read permission or is note owner
	if readPermissions(r) {

		//Get username
		username := getUserName(r)
		noteid := getID(r)
		note = getPartOfNoteSQL(noteid, username)

	}
	return note
}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR SQL TO DATABASE
// ------------------------------------------------------------------------------------------------------------------------------------------------

func getOwnedNote(r *http.Request) (note Note) {
	if noteOwner(r) {

		//get username
		username := getUserName(r)
		noteid := getID(r)
		note = getOwnedNoteSQL(noteid, username)

	}
	return note
}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR SQL TO DATABASE
// ------------------------------------------------------------------------------------------------------------------------------------------------

//Update note that you own
func updateOwnedNote(r *http.Request) (success bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()
	//get username
	username := getUserName(r)
	//get the note id that we need to update
	noteid := getID(r)
	//get values from form
	title := r.FormValue("title")
	body := r.FormValue("body")

	//validate that there is text in title
	if validateInput(title) {

		if updateOwnedNoteSQL(title, body, noteid, username) {
			return true
		}
		return false

	}

	return false

}

// ------------------------------------------------------------------------------------------------------------------------------------------------
// CLEAR SQL TO DATABASE
// ------------------------------------------------------------------------------------------------------------------------------------------------

//Update note you are part of
func updatePartOfNote(r *http.Request) bool {
	if userStillLoggedIn(r) {

		//get username
		//username := getUserName(r)
		//get note id
		noteID := getID(r)

		//Get value from form
		body := r.FormValue("body")

		if updatePartOfNoteSQL(noteID, body) {
			return true
		}
	}
	return false
}

//Need to move to database.go but just putting it here for use sake
//Add more users to a note
func addPermissions(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {

		noteid := getID(r)

		var read string
		var write string

		users := r.Form["user"]
		for _, user := range users {
			//get included checkbox value
			includedCheckbox := r.FormValue("includedCheckbox_" + user)
			//Check that the user has been included
			if includedCheckbox != "" {
				read = "t"
				writeCheckbox := r.FormValue("writeCheckbox_" + user)
				//Check that the user has write privlages
				if writeCheckbox != "" {
					write = "t"
				} else {
					write = "f"
				}
				//Add to database
				if !addPermissionSQL(noteid, user, read, write) {
					http.Error(w, "Database Error", http.StatusInternalServerError)
				}

			}

		}

		http.Redirect(w, r, "/listNotes", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func editPrivileges(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) && noteOwner(r) {
		var included string
		var writeValue string
		var write string
		noteid := getID(r)
		users := r.Form["user"]
		for _, user := range users {
			included = r.FormValue("includedCheckbox_" + user)
			//If user not included anymore note privilege will be moved
			if included == "" {
				removePrivilege(noteid, user)
				continue
			}
			writeValue = r.FormValue("writeCheckbox_" + user)
			if writeValue != "" {
				write = "t"
			} else {
				write = "f"
			}
			updatePrivilege(noteid, user, write)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//get note id from url
func getID(r *http.Request) (id string) {
	id = mux.Vars(r)["id"]
	return id
}

////////////GROUP PRIVILEGES SECTION/////////////////////

//Displays HTML page for showing all groups that the user owns
func viewAllSavedGroups(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		username := getUserName(r)
		//get all groups belonging to the user
		groups := getAllGroups(username)

		tpl.ExecuteTemplate(w, "viewGroups.gohtml", groups)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//displays the html page for creating a new group
func createGroup(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		username := getUserName(r)
		users := listAllUsersSQL(username)
		tpl.ExecuteTemplate(w, "createGroup.gohtml", users)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Add a new Group to the database
func addGroup(w http.ResponseWriter, r *http.Request) {
	//Make sure user is still logged in
	if userStillLoggedIn(r) {

		var newGroup Group
		alert := ""
		newGroup.GroupOwner = getUserName(r)
		newGroup.GroupTitle = r.FormValue("title")

		//Validate that user has entered a title
		if validateInput(newGroup.GroupTitle) {
			writePrivilege := r.FormValue("write_privilege")
			//Validate that group has write privileges
			if writePrivilege != "" {
				newGroup.GroupWrite = "t"
			} else {
				newGroup.GroupWrite = "f"
			}

			//Create new group
			groupid := createNewGroup(newGroup.GroupTitle, newGroup.GroupOwner, "t", newGroup.GroupWrite)
			//validate that group was added to database
			if groupid == "" {
				alert = `<script>
					alert("Group Not Successfully Created");
					window.location.href="/";
				</script>`
			} else {
				alert = `<script>
					alert("Successfully Created Group");
					window.location.href="/";
				</script>`
			}
			users := r.Form["user"]

			for _, user := range users {
				includedCheckBox := r.FormValue("includedCheckbox_" + user)
				//validate that user is included in the group
				if includedCheckBox != "" {
					saveGroupUserSQL(groupid, user)
				}
			}
		} else {
			alert = `<script>
					alert("Please enter a title");
					window.location.href="/createGroup";
				</script>`
		}
		w.Write([]byte(alert))
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Delete Group
func deleteGroup(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		username := getUserName(r)
		alert := ""
		groupid := getID(r)
		//Validate that group owner
		if validateGroupOwner(username, groupid) {

			//Delete note and validate that note has been deleted
			if removeGroup(groupid) {
				alert = `<script>
							alert("Successfully Deleted Group");
							window.location.href="/";
							</script>`
			} else {
				alert = `<script>
							alert("Group was not successfully deleted");
							window.location.href="/";
							</script>`
			}

			w.Write([]byte(alert))
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Display the edit group page
func viewEditGroupUsers(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		groupid := getID(r)
		username := getUserName(r)
		//validate that user is note owner
		if validateGroupOwner(username, groupid) {

			users := getGroupUsers(groupid)

			//Make temp struct to hold group id and users
			editStruct := struct {
				GroupID string
				Users   []string
			}{
				GroupID: groupid,
				Users:   users,
			}
			tpl.ExecuteTemplate(w, "editGroupUsers.gohtml", editStruct)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//edit group users
func editGroupUsers(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		success := true
		alert := ""
		groupid := getID(r)
		username := getUserName(r)
		//For weird bug DO NOT REMOVE
		_ = r.FormValue("includedCheckbox_test")
		//validate that user is the owner of the group
		if validateGroupOwner(username, groupid) {
			users := r.Form["user"]
			for _, user := range users {
				included := r.FormValue("includedCheckbox_" + user)
				//Validate that user has been removed
				if included == "" {
					if !removeGroupUser(groupid, user) {
						success = false
						break
					}
				}
			}
			if success {
				alert = `<script>
							alert("Successfully Edited Group");
							window.location.href="/";
							</script>`
			} else {
				alert = `<script>
							alert("Group was not successfully Edited");
							window.location.href="/";
						</script>`
			}
			//notify user
			w.Write([]byte(alert))
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

//Display the Group the user wants to edit
func viewGroup(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		username := getUserName(r)
		groupid := getID(r)
		//validate the the user is the owner of the note
		if validateGroupOwner(username, groupid) {
			group := getGroup(groupid)

			tpl.ExecuteTemplate(w, "editGroup.gohtml", group)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//update a group
func updateGroup(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		username := getUserName(r)
		groupid := getID(r)
		var write string
		var alert string
		//Validate that user is the owner of the group
		if validateGroupOwner(username, groupid) {
			writeCheckbox := r.FormValue("writeCheckbox")
			//Validate if group can have write permissions
			if writeCheckbox == "" {
				write = "f"
			} else {
				write = "t"
			}
			//Update group
			if editGroup(groupid, write) {
				alert = `<script>
							alert("Successfully Updated Group");
							window.location.href="/";
						</script>`
			} else {
				alert = `<script>
							alert("Group was not successfully Updated");
							window.location.href="/";
						</script>`
			}
			//Show user result
			w.Write([]byte(alert))
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Display users that the user wants to add to a group
func viewAddGroupUsers(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		username := getUserName(r)
		groupid := getID(r)
		//Validate that user is the owner of the group
		if validateGroupOwner(username, groupid) {

			users := getAvaliableGroupUsers(groupid, username)
			//Temp struct
			addUsersStruct := struct {
				GroupID string
				Users   []string
			}{
				GroupID: groupid,
				Users:   users,
			}

			tpl.ExecuteTemplate(w, "addGroupUsers.gohtml", addUsersStruct)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Add the users to the group
func addGroupUsers(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		username := getUserName(r)
		groupid := getID(r)
		var alert string
		success := true
		//Valdiate if user is the owner of the note
		if validateGroupOwner(username, groupid) {
			//for dealing with bug DO NOT REMOVE
			_ = r.FormValue("includedCheckbox_test")
			users := r.Form["user"]
			for _, user := range users {
				includedCheckBox := r.FormValue("includedCheckbox_" + user)
				//Validate that add user has been included
				if includedCheckBox != "" {
					//Validate that add user has been added to database
					if !saveGroupUserSQL(groupid, user) {
						success = false
					}
				}
			}
			//Make sure that the users were added successfully
			if success {
				alert = `<script>
							alert("Successfully Added Users to Group");
							window.location.href="/";
						</script>`
			} else {
				alert = `<script>
							alert("Users were not successfully Added to Group");
							window.location.href="/";
						</script>`
			}

			//Display result to user
			w.Write([]byte(alert))
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

// Validate that user has entered value
func validateInput(input string) bool {
	if input == "" {
		return false
	}

	return true
}

////////////////ANALYSE SECTION//////////////

//Show a note that the user wants to analyse
func showAnalyseNote(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		var note Note

		note = getNoteToAnalyse(r)

		tpl.ExecuteTemplate(w, "analyseNote.gohtml", note)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Analyse a note that the user has choosen
func analyseNote(w http.ResponseWriter, r *http.Request) {
	if userStillLoggedIn(r) {
		//Get noote to be analises
		note = getNoteToAnalyse(r)
		//get pattern from form
		pattern := r.FormValue("pattern")
		//See how many times pattern is in the note body
		count := strconv.Itoa(strings.Count(note.NoteBody, pattern))

		//create prompt
		alert := `<script>
					alert("` + pattern + ` is in ` + note.NoteTitle + ` ` + count + ` times");
					window.location.href="/analyseNote/` + strconv.Itoa(note.NoteID) + `";
					</script>`
		w.Write([]byte(alert))
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Get note to analyse
func getNoteToAnalyse(r *http.Request) (note Note) {
	if noteOwner(r) {
		note = getOwnedNote(r)
	} else {
		note = getPartOfNote(r)
	}

	return note
}
