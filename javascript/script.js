function deleteNote(noteID) {
    if (window.confirm("By clicking Okay the note will be deleted")) {
        window.location.href = "/deleteNote/" + noteID;
    }
}

function saveUserModal() {
    console.log("Entered save");
    $('#addUserModal').modal('hide');

}
