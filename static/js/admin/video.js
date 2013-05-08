$(function () {
    createSortable();
    $('a.addVideo').click(function (event) {
        event.preventDefault();
        var idstr = $(this).attr('id').split(':')[1];
        $.getJSON("/Video/AddVideo", { 'ytID': idstr }, function (data) {
            $("#liveVideos").sortable("destroy");
            $('#liveVideos').append('<li class="sortableVideo" id="video_' + data.ID + '"><img src="' + data.Screenshot + '" alt="' + data.Title + '" /><span class="videotitle">' + data.Title + '</span><br /><a class="deleteVideo" href="#" id="delete_' + data.ID + '">Remove</a><span class="clear"></span></li>');
            createSortable();
        });
    });
    $(document).on('click','a.deleteVideo', function (event) {
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