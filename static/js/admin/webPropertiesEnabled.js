$(document).ready(function () {

    var dt = $('table').dataTable({
        "bJQueryUI": true
    });


    $('.isEnabled').live('click', function () {
        var record_id = $(this).attr('id').split(':')[1];
        set_isEnabled(record_id);
    });

    $('.isFinalApproved').live('click', function () {
        var record_id = $(this).attr('id').split(':')[1];
        set_isFinalApproved(record_id);
    });

    $('.isDenied').live('click', function () {
        var record_id = $(this).attr('id').split(':')[1];
        set_isDenied(record_id);
    });

    $('.searchByColor').click(function () {
        dt.fnFilter($(this).attr('title'));
        $('.dataTables_filter').find("input[type=text]").first().attr('value', $(this).attr('title'));
    });

});


/*
* This function is going to make an AJAX call to the controller and set the isEnabled field.
* @param record_id: Primary Key for web property
*/
function set_isEnabled(record_id) {
    $.get('/Customers/SetWebPropertyStatus', { 'record_id': record_id }, function (response) {
        if (response != '') {
            showMessage(response);
        } else {
            showMessage("The Web Property's approved pending status has been changed. Refresh page to see pending date.");
        }
    },"html");
}


function set_isFinalApproved(record_id) {
    $.get('/Customers/WPSetIsFinalApproved', { 'record_id': record_id }, function (response) {
        if (response != '') {
            showMessage(response);
        } else {
            showMessage("The Web Property has been officially approved.");
        }
    }, "html");
}


function set_isDenied(record_id) {
    $.get('/Customers/WPSetIsDenied', { 'record_id': record_id }, function (response) {
        if (response != '') {
            showMessage(response);
        } else {
            showMessage("The Web Property has been rejected.");
        }
    }, "html");
}




