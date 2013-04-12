$(document).ready(function () {

    var group_table = $('table').dataTable({ 'bJQueryUI': true });

    $('.delete').live('click', function () {
        if (confirm('Are you sure you want to remove this Forum Group?')) {
            var path = $(this).attr('href');
            window.location = path;
        }
        return false;
    });

    $('.deletetopic').live('click', function () {
        if (confirm('Are you sure you want to remove this Topic?')) {
            var path = $(this).attr('href');
            window.location = path;
        }
        return false;
    });

});