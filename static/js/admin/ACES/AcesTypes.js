$(function () {
    AcesTable = $('#AcesTypes').dataTable({
        "bJQueryUI": true
    });
    $(document).on('click', '.removeACESType', function (e) {
        e.preventDefault();
        var idstr = $(this).data('id');
        var table_row = $(this).parent().parent().get()[0];
        if (confirm('Are you sure you want to remove this ACES Type?')) {
            $.get('/ACES/RemoveACESType/' + idstr, function (response) {
                if (response.length != 0) {
                    AcesTable.fnDeleteRow(table_row);
                    showMessage("ACES Type deleted successfully.");
                }
            });
        }
    });
});