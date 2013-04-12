$(function () {

    var rep_table = $('table').dataTable({ 'bJQueryUI': true });

    $('.delete').live('click', function () {
        if (confirm('Are you sure you want to remove this sales rep?')) {
            var idstr = $(this).data('id');
            var table_row = $(this).parent().parent().get()[0];
            $.get("/SalesRep/Delete/" + idstr, function (resp) {
                if (resp.length == 0) {
                    rep_table.fnDeleteRow(table_row);
                    showMessage('Sales Rep Removed.');
                } else {
                    showMessage('There was a problem');
                }
            });
        }
        return false;
    });
});