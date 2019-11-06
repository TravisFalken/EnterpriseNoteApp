function deleteNote(noteID) {
    if (window.confirm("By clicking Okay the note will be deleted")) {
        window.location.href = "/deleteNote/" + noteID;
    }
}

function saveUserModal() {
    console.log("Entered save");
    $('#addUserModal').modal('hide');

}

//Something weird
function editNotePage(noteID){
    if(noteID > 0){
        window.location.href = "/editNote/" +noteID;
    }
   
}

function updateNote(noteID){
    //var form = $('#editForm');
    var form = $('#editForm');
    //var formData = new FormData(form);
        if(noteID > 0){
            window.alert("Note Has been updated!");
            
            $.ajax({
                url: '/updateNote/' + noteID,
                type: 'POST',
                processData: false,
                contentType: false,
                date: form.serialize()
            });
        }
        
    
}

//Want to add users to note
function addUsersClicked(noteID){
    window.location.href = "/addUsers/" + noteID;
    /*
    $.ajax({
        url: '/addUsers/' + noteID,
        type: 'GET'
    });
    */
}

//Want to list all of the privileges for a note
function editPrivileges(noteID){
    window.location.href = "/listPrivileges/" + noteID;
}

//Delete group
function deleteGroup(groupID){
    if (window.confirm("By clicking Okay the Group will be deleted")) {
        window.location.href = "/deleteGroup/" + groupID;
    }
    
}

//Edit group users
function editGroupUsers(groupID){
    window.location.href = "/viewEditGroupUsers/" + groupID;
}

//view a group to edit
function viewGroup(groupID){
    window.location.href = "/viewGroup/" + groupID;
}

//view users to add to note
function viewaddGroupUsers(groupID){
    window.location.href = "/AddUsersGroup/" + groupID;
}

