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

