$(function () {

    $('#choose-photo').click(function () { chooseFile(); })
    $('#clear-photo').click(function () {
        if (confirm('Are you sure you want to clear this image?')) {
            $('#photo-file img').attr('src', '/Content/img/noimage.jpg').attr('alt', 'No Photo');
            $('#file').val('');
        }
    })

});

function selectFile(url) {
    $('#file').val(url);
    $('#photo-file img').attr('src', url).attr('alt', 'Photo');
    $("#file-dialog").dialog("close");
    $("#file-dialog").empty();
}