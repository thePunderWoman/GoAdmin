$(document).ready(function () {

    var group_table = $('table').dataTable({
        "bJQueryUI": true,
        "aoColumns": [
            null,
            null,
            { "sType": "date" },
            null,
            null
        ],
        "aaSorting": [[2, "asc"]]
    });

    $('.delete').live('click', function () {
        if (confirm('Are you sure you want to remove this Thread? \nDon\'t worry! This CAN be undone.')) {
            var path = $(this).attr('href');
            window.location = path;
        }
        return false;
    });
});