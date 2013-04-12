$(document).ready(function () {

    var dt = $('#webRequirements').dataTable({
        "bJQueryUI": true,
        "iDisplayLength": 50,
        "aaSorting": []
    });

    $('.reqCheck').live('click', function () {
        var webPropRequirementID = $(this).attr('id').split(':')[1];// web prop requirement
        var webPropID = $(this).attr('id').split(':')[2]; // web prop id
        set_reqCheck(webPropRequirementID, webPropID);
    });

});

function set_reqCheck(webPropRequirementID, webPropID) {
    $.get('/Customers/SetWebPropertyRequirementStatus', { 'webPropRequirementID': webPropRequirementID, 'webPropID': webPropID }, function (response) {
        if (response != '') {
            showMessage(response);
        } else {
            showMessage("Authorized Dealer Requirements have been updated.");
        }
    }, "html");
}
