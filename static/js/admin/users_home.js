$(document).ready(function () {

    $('table').dataTable({
        "bJQueryUI": true
    });

    $('#choose-photo').click(function () { chooseFile(); })
    $('#clear-photo').click(function () {
        if (confirm('Are you sure you want to clear this image?')) {
            $('#photo-file img').attr('src', '/Content/img/nophoto.jpg').attr('alt', 'No Photo');
            $('#file').val('');
        }
    })

    $('.user_action').live('change', function () {
        var action = $(this).val();
        var user_id = $(this).attr('id').split(':')[1];
        switch (action) {
            case 'view':
                $.get('/Users/GetUser', { 'user_id': user_id }, quickView, 'json');
                break;
            case 'edit':
                window.location.href = "/Users/Edit?user_id=" + user_id;
                break;
            case 'delete':
                if (confirm('Are you sure you want to remove this user?')) {
                    $.get('/Users/RemoveUser', { 'userID': user_id }, function (data) {
                        deleteUser(data, user_id);
                    });
                }

                break;
            default:

        }
        $(this).val(0);
    });

    $('.isActive').live('click', function () {
        var user_id = $(this).attr('id').split(':')[1];
        set_isActive(user_id);
    });

});


/*
* This function is going to make an AJAX call to the controller and set the isActive field.
* @param userID: Primary Key for user
*/
function set_isActive(userID) {
    $.get('/Users/SetUserStatus',{'userID':userID},function(response){
        if (response != '') {
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
    if (response != '') {
        showMessage('There was an error while removing the user.');
    } else {
        $('#user\\:' + user_id).remove();
        showMessage('User was successfully removed.');
    }
}

/*
* Does an ajax call to get the data for this user and display it on this page.
* @param userID: Primary Key for user
*/
function quickView(user) {
    if (user != '') {
        $('#user_name').find('h4').text(user.fname + ' ' + user.lname);
        $('#username').text(user.username);
        $('#email').text(user.email);
        $('#website').text(user.website);
        $('#phone').text(user.phone);
        $('#fax').text(user.fax);
        $('#user_container').slideDown();
    }
}

function selectFile(url) {
    $('#file').val(url);
    $('#photo-file img').attr('src', url).attr('alt', 'Photo');
    $("#file-dialog").dialog("close");
    $("#file-dialog").empty();
}