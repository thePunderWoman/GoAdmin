var videoTable = "";
$(document).ready(function () {
    videoTable = $('table').dataTable({ "bJQueryUI": true });

    $('#addVideo').live('click', function () {
        clearFormValues();
        showVideoForm(0, '', '', false);
    });

    $(document).on('click', ".edit", function (event) {
        event.preventDefault();
        var vidID = $(this).data("id");
        videoTable.fnDeleteRow($(this).parent().parent().get()[0]);
        $.getJSON("/Product/GetPartVideo", { 'pVideoID': vidID }, function (data) {
            clearFormValues();
            showVideoForm(data.pVideoID, data.vTypeID, data.video, data.isPrimary);
        });
    });

    $(document).on('click', ".remove", function (event) {
        event.preventDefault();
        var vidID = $(this).data('id');
        var removelink = $(this);
        if (confirm('Are you sure you want to remove this video?')) {
            $.post('/Product/DeleteVideo', { 'videoID': vidID }, function (response) {
                if (response == "success") {
                    videoTable.fnDeleteRow($(removelink).parent().parent().get()[0]);
                    showMessage("Video Removed Successfully");
                } else {
                    showMessage("Error Removing Video");
                }
            }, 'text');
        }
    });

    $(document).on('click', '#btnReset', function () {
        var vidID = Number($('#pVideoID').val());
        if (vidID > 0) {
            $.getJSON("/Product/GetPartVideo", { 'pVideoID': vidID }, function (response) {
                videoTable.fnAddData([
                            response.videoType.name,
                            '<iframe data-video="' + response.video + '" width="177" height="140" src="http://www.youtube.com/embed/' + response.video + '" frameborder="0" allowfullscreen></iframe>',
                            ((response.isPrimary) ? 'Yes' : 'No'),
                            '<a href="javascript:void(0)" class="edit" data-id="' + response.pVideoID + '" title="Edit Video">Edit</a> | <a href="javascript:void(0)" class="remove" data-id="' + response.pVideoID + '" title="Remove Video">Remove</a>'
                            ]);
            });
        }
        clearFormValues();
    });

    $(document).on('click', '#btnSave', function (event) {
        event.preventDefault();
        var type = $('#videoType').val();
        var video = $('#video').val();
        var vidID = $('#pVideoID').val();
        var partID = $('#partID').val();
        var isPrimary = $('#isPrimary').is(':checked');
        if (video == "" || type == "") {
            showMessage("You must enter a video ID and select a video Type.");
            return;
        }
        $.post("/Product/SaveVideo", { 'partID': partID, 'pVideoID': vidID, 'video': video, 'videoType': type, 'isPrimary': isPrimary }, function (response) {
            if (response.error == null) { // Success 0FLQ4rACE-0
                videoTable.fnAddData([
                        response.videoType.name,
                        '<iframe data-video="' + response.video + '" width="177" height="140" src="http://www.youtube.com/embed/' + response.video + '" frameborder="0" allowfullscreen></iframe>',
                        ((response.isPrimary) ? 'Yes' : 'No'),
                        '<a href="javascript:void(0)" class="edit" data-id="' + response.pVideoID + '">Edit</a> | <a href="javascript:void(0)" class="delete" data-id="' + response.pVideoID + '">Delete</a>'
                        ]);
                clearFormValues();
            } else {
                showMessage(response.error);
            }
        }, 'json');
    });
});

function showVideoForm(pVideoID, videoType, video, isPrimary) {
    $('#pVideoID').val(pVideoID);
    $('#videoType').val(videoType);
    $('#video').val(video);
    $('#isPrimary').attr('checked', false);
    if (isPrimary) {
        $('#isPrimary').attr('checked', true);
    }
    $('.form_left').slideDown();
}

function clearFormValues() {
    $('#pVideoID').val(0);
    $('#isPrimary').attr('checked', false);
    $('#video').val('');
    $('#videoType').val('');
    $('.form_left').slideUp();
}