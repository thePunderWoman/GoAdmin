$(function () {

    $('table').dataTable({
        "bJQueryUI": true
    });

    $('#choose-photo').on("click",function () { chooseFile(); })
    $('#clear-photo').on("click",function () {
        if (confirm('Are you sure you want to clear this image?')) {
            $('#photo-file img').attr('src', '/Content/img/nophoto.jpg').attr('alt', 'No Photo');
            $('#file').val('');
        }
    })
    $(document).on('change','.user_action', function () {
        var action = $(this).val();
        var user_id = $(this).attr('id').split(':')[1];
        switch (action) {
            case 'edit':
                window.location.href = "/Users/Edit/" + user_id;
                break;
            case 'delete':
                if (confirm('Are you sure you want to remove this user?')) {
                    $.getJSON('/Users/Delete/' + user_id, function (data) {
                        deleteUser(data, user_id);
                    });
                }

                break;
            default:

        }
        $(this).val(0);
    });

    $('.isActive').on('click', function () {
        var user_id = $(this).attr('id').split(':')[1];
        set_isActive(user_id);
    });

});


/*
* This function is going to make an AJAX call to the controller and set the isActive field.
* @param userID: Primary Key for user
*/
function set_isActive(userID) {
    $.getJSON('/Users/SetUserStatus/' + userID,function(response){
        if (!response.Success) {
            showMessage(response);
        }else{
            showMessage("User's status has been updated.");
        }
    });
}


/*
* Makes a call to the controller and removes the user from the database.
* @param userID: Primary Key for user
*/
function deleteUser(response, user_id) {
    if (!response.Success) {
        showMessage('There was an error while removing the user.');
    } else {
        $('#user\\:' + user_id).remove();
        showMessage('User was successfully removed.');
    }
}

function selectFile(url) {
    $('#file').val(url);
    $('#photo-file img').attr('src', url).attr('alt', 'Photo');
    $("#file-dialog").dialog("close");
    $("#file-dialog").empty();
}