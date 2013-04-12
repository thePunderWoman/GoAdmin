$(function () {
    AcesTable = $('#ConfigAttributes').dataTable({
        "bJQueryUI": true
    });
    $(document).on('click', '.remove', function (e) {
        e.preventDefault();
        var idstr = $(this).data('id');
        var table_row = $(this).parent().parent().get()[0];
        if (confirm('Are you sure you want to remove this Configuration Attribute?')) {
            $.get('/ACES/RemoveConfigurationAttribute/' + idstr, function (response) {
                if (response.length != 0) {
                    AcesTable.fnDeleteRow(table_row);
                    showMessage("Configuration Attribute deleted successfully.");
                }
            });
        }
    });
});