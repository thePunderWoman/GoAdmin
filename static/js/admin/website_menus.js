$(document).ready(function () {

    var menutable = $('table').dataTable({
        "bJQueryUI": true
    });

    $('.menu_action').live('change', function () {
        var action = $(this).val();
        var menu_id = $(this).attr('id').split(':')[1];
        var table_row = $(this).parent().parent().get()[0];
        switch (action) {
            case 'manage':
                window.location.href = "/Website/Content/Menu/" + menu_id;
                break;
            case 'edit':
                window.location.href = "/Website/EditMenu/" + menu_id;
                break;
            case 'delete':
                if (confirm('Are you sure you want to remove this menu?')) {
                    $.post('/Website/RemoveMenu', { 'menuID': menu_id }, function (data) {
                        menutable.fnDeleteRow(table_row);
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
function deleteMenu(response, menu_id) {
    if (response != '') {
        showMessage('There was an error while removing the menu.');
    } else {
        $('#menu\\:' + menu_id).remove();
        showMessage('Menu was successfully removed.');
    }
}
