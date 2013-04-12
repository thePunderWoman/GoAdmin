$(document).ready(function () {

    var table = $('#resource_table').dataTable({
        "bJQueryUI": true
    });

    var user_table = $('#user_table').dataTable({
        "bJQueryUI": true
    });

    $('#addNew').live('click', function () {
        loadForm(0, '', '', '', '', '');
    });

    $('#btnReset').live('click', function () {
        clearForm();
    });

    $('#btnSubmit').live('click', function () {
        var resourceID = $('#resourceID').val();
        var name = $('#resource_name').val();
        var url = $('#resource_url').val();
        var username = $('#resource_username').val();
        var pass = $('#resource_password').val();
        var comments = $('#comments').val();

        if (resourceID == 0) { // Add new
            if (name.length > 0) {
                if (confirm("Are you sure you want to save this resource?")) {
                    $.post("/Account/Add", { "name": name, "url": url, "username": username, "password": pass, "comments": comments }, function (response) {
                        if (response.error == null) {
                            table.fnAddData([
                                    response.resource_name,
                                    response.resource_url,
                                    response.username,
                                    response.password,
                                    response.comments,
                                    '<select class="action" id="' + response.resourceID + '"><option value="">- Select Action -</option><option value="edit">Edit</option><option value="users">Manage Users</option><option value="delete">Remove</option></select>'
                                ]);
                            clearForm();
                        } else {
                            showMessage(response.error);
                        }
                    }, 'json');
                }
            } else {
                showMessage("You must enter a resource name.");
            }
        } else { // Update
            if (name.length > 0) {
                if (confirm("Are you sure you want to save this resource?")) {
                    table.fnDeleteRow($('#' + resourceID).parent().parent().get()[0]);
                    $.post("/Account/Update", { "resourceID": resourceID, "name": name, "url": url, "username": username, "password": pass, "comments": comments }, function (response) {
                        if (response.error == null) {
                            table.fnAddData([
                                    response.resource_name,
                                    response.resource_url,
                                    response.username,
                                    response.password,
                                    response.comments,
                                    '<select class="action" id="' + response.resourceID + '"><option value="">- Select Action -</option><option value="edit">Edit</option><option value="users">Manage Users</option><option value="delete">Remove</option></select>'
                                ]);
                            clearForm();
                        } else {
                            showMessage(response.error);
                        }
                    }, 'json');
                }
            } else {
                showMessage("You must enter a resource name.");
            }
        }
        return false;
    });

    $('.action').live('change', function () {
        var resourceID = $(this).attr('id');
        var action = $(this).val();
        var name = $(this).parent().prev().prev().prev().prev().prev().text();
        var url = $(this).parent().prev().prev().prev().prev().text();
        var username = $(this).parent().prev().prev().prev().text();
        var password = $(this).parent().prev().prev().text();
        var comments = $(this).parent().prev().text();
        switch (action) {
            case 'edit':
                loadForm(resourceID, name, url, username, password, comments);
                break;
            case 'users':
                clearForm();
                // Get the users that can view this resource
                $.getJSON('/Account/GetResourceUsers', { 'resourceID': resourceID }, function (users) {
                    $.each(users, function (i, user) {
                        user_table.fnAddData([
                                user.user,
                                user.username,
                                '<a href="javascript:void(0)" class="removeUser" id="' + user.userID + '">Remove</a>'
                            ]);
                    });
                });
                $('#user_resourceID').val(resourceID);
                $('#user_form').slideDown();

                break;
            case 'delete':
                if (confirm("Are you sure you want to remove this resource?")) {
                    $.post("/Account/Remove", { "resourceID": resourceID }, function (response) {
                        if (response.length == 0) { // Remove table row
                            table.fnDeleteRow($('#' + resourceID).parent().parent().get()[0]);
                            clearUserForm();
                        } else {
                            showMessage(response);
                        }
                    });
                }
                break;
            default:
                // DO NOTHING
        }
        $(this).val('');
    });

    $('#btnAddUser').live('click', function () {
        var resourceID = $('#user_resourceID').val();
        var userID = $('#newUser').val();
        $.post('/Account/AddUserToResource', { 'resourceID': resourceID, 'userID': userID }, function (response) {
            if (response.error == null) {
                user_table.fnAddData([
                    response.user,
                    response.username,
                    '<a href="javascript:void(0)" class="removeUser" id="' + response.userID + '">Remove</a>'
                    ]);
                $('#newUser').val(0);
            } else {
                showMessage(response.error);
            }
        }, 'json');

    });

    $('#btnClose').live('click', function () {
        clearUserForm();
    });

    $('.removeUser').live('click', function () {
        var resourceID = $('#user_resourceID').val();
        var userID = $(this).attr('id');
        var table_row = $(this).closest('tr').get()[0];
        if (resourceID > 0 && userID > 0) {
            $.post('/Account/RemoveUserFromResource', { 'resourceID': resourceID, 'userID': userID }, function (response) {
                if (response == null) {
                    user_table.fnDeleteRow(table_row);
                } else {
                    showMessage(response);
                }
            }, 'json');
        } else {
            showMessage("Invalid data.");
        }
    });

});

function clearUserForm() {
    user_table.fnClearTable();
    $('#user_resourceID').val(0);
    $('#user_form').slideUp();
}

function clearForm() {
    $('#resourceID').val(0);
    $('#resource_name').val('');
    $('#resource_url').val('');
    $('#resource_username').val('');
    $('#resource_password').val('');
    $('#comments').val('');
    $('.form_left').slideUp();
    $('#addNew').fadeIn();
}

function loadForm(resourceID, name, url, username, password, comments) {
    $('#user_form').slideUp();
    $('#resourceID').val(resourceID);
    $('#resource_name').val(name);
    $('#resource_url').val(url);
    $('#resource_username').val(username);
    $('#resource_password').val(password);
    $('#comments').val(comments);
    $('.form_left').slideDown();
    $('#addNew').fadeOut();
}