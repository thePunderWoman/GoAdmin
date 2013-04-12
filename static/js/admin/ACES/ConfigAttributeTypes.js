$(function () {
    AcesTable = $('#ConfigAttributeTypes').dataTable({
        "bJQueryUI": true
    });
    $(document).on('click', '.remove', function (e) {
        e.preventDefault();
        var idstr = $(this).data('id');
        var table_row = $(this).parent().parent().get()[0];
        if (confirm('Are you sure you want to remove this Configuration Type?')) {
            $.get('/ACES/RemoveConfigurationType/' + idstr, function (response) {
                if (response.length != 0) {
                    AcesTable.fnDeleteRow(table_row);
                    showMessage("Configuration Type deleted successfully.");
                }
            });
        }
    });
});