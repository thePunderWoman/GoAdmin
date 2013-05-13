$(function () {

    var rep_table = $('table').dataTable({ 'bJQueryUI': true });

    $(document).on('click','.delete', function () {
        if (confirm('Are you sure you want to remove this sales rep?')) {
            var idstr = $(this).data('id');
            var table_row = $(this).parent().parent().get()[0];
            $.getJSON("/SalesRep/Delete/" + idstr, function (resp) {
                if (resp) {
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