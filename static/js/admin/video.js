$(function () {
    createSortable();
    $('a.addVideo').click(function (event) {
        event.preventDefault();
        var idstr = $(this).attr('id').split(':')[1];
        $.getJSON("/Video/AddVideo", { 'ytID': idstr }, function (data) {
            $("#liveVideos").sortable("destroy");
            $('#liveVideos').append('<li class="sortableVideo" id="video_' + data.videoID + '"><img src="' + data.thumb + '" alt="' + data.videoTitle + '" /><span class="videotitle">' + data.videoTitle + '</span><br /><a class="deleteVideo" href="#" id="delete_' + data.videoID + '">Remove</a><span class="clear"></span></li>');
            createSortable();
        });
    });
    $('a.deleteVideo').live("click", function (event) {
        event.preventDefault();
        if (confirm("Are you sure you want to remove this video?")) {
            var idstr = $(this).attr('id').split('_')[1];
            $.post("/Video/Delete", { "id": idstr }, function (data) {
                if (data == "success") {
                    $('#video_' + idstr).remove();
                }
            }, "text");
        }
    });
});

function updateSort() {
    var x = $('#liveVideos').sortable("serialize");
    $.post("/Video/updateSort?" + x);
}

function createSortable() {
    $("#liveVideos").sortable({
        handle: 'img',
        update: function (event, ui) { updateSort(); }
    }).disableSelection();
}