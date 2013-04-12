$(function () {

    var table = $('table').dataTable({ 'bJQueryUI': true });

    $('.delete').live('click', function () {
        if (confirm('Are you sure you want to restore access to this IP?')) {
            var path = $(this).attr('href');
            window.location = path;
        }
        return false;
    });

});