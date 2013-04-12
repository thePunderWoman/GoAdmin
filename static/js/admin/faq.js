$(document).ready(function () {

    var faq_table = $('table').dataTable({ 'bJQueryUI': true });

    $('.delete').live('click', function () {
        if (confirm('Are you sure you want to remove this question and answer?')) {
            var path = $(this).attr('href');
            var table_row = $(this).parent().parent().get()[0];
            $.get(path, function (resp) {
                if (resp.length == 0) {
                    faq_table.fnDeleteRow(table_row);
                    showMessage('FAQ Removed.');
                } else {
                    showMessage(resp);
                }
            });
        }
        return false;
    });

    $('#btnSave').click(function () {
        var errs = 0;
        var q = $.trim($('#question').val());
        if (q.length == 0) {
            $('#question').after('<br /><p style="color:#bf0000;margin-top:0px">You must enter a question.</p>');
            $('#question').addClass('required');
            errs++;
        }

        var a = $.trim($('#answer').val());
        if (a.length == 0) {
            $('#answer').after('<br /><p style="color:#bf0000;margin-top:0px">You must enter an answer.</p>');
            $('#answer').addClass('required');
            errs++;
        }

        if (errs > 0) {
            showMessage('Please fix the errors.');
            return false;
        }
    });
});