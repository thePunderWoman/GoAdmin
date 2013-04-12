$(function () {
    $('#revisions table').dataTable({ "bJQueryUI": true });
    CKEDITOR.replace('page_content', {
        filebrowserImageUploadUrl: '/File/CKUpload',
        filebrowserImageBrowseUrl: '/File/CKIndex',
        filebrowserImageWindowWidth: '640',
        filebrowserImageWindowHeight: '480'
    });

    $('#btnSubmitPublish').click(function () {
        $('#publish').val("true");
        $('#saveform').submit();
    });

    $('tr.active td a.delete').hide();
    $('tr.active td span.delbefore').hide();

    var revID = $('#revID').val();
    $('#rev' + revID + ' td a.edit').hide();
    $('#rev' + revID + ' span.editafter').hide();

    $(document).on('click', '.delete', function (e) {
        e.preventDefault();
        var url = $(this).attr('href');
        if(confirm('Are you sure you want to delete this revision?')) {
            window.location.href = url;
        }
    });
    $(document).on('click', '.activate', function (e) {
        e.preventDefault();
        var url = $(this).attr('href');
        if (confirm('Are you sure you want to make this revision active?')) {
            window.location.href = url;
        }
    });
});
