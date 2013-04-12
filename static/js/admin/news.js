$(document).ready(function () {

    var news_table = $('table').dataTable({ 'bJQueryUI': true });

    $('.delete').live('click', function () {
        if (confirm('Are you sure you want to remove this news item?')) {
            var path = $(this).attr('href');
            var table_row = $(this).parent().parent().get()[0];
            $.get(path, function (resp) {
                if (resp.length == 0) {
                    news_table.fnDeleteRow(table_row);
                    showMessage('News Item Removed.');
                } else {
                    showMessage(resp);
                }
            });
        }
        return false;
    });

    $('#btnSave').click(function () {
        var errs = 0;
        var t = $.trim($('#title').val());
        if (t.length == 0) {
            $('#title').after('<br /><p style="color:#bf0000;margin-top:0px">You must enter a title.</p>');
            $('#title').addClass('required');
            errs++;
        }

        if (errs > 0) {
            showMessage('Please fix the errors.');
            return false;
        }
    });
});