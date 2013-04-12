$(function () {
    $('div.controls a').click(function (event) {
        event.preventDefault();
    });

    $('.edit').click(function () {
        var id = $(this).data("id");
        $.getJSON('/Forum/GetPost', { 'postid': id }, function (data) {
            $('#postID').val(data.postID);
            $('#parentid').val(0);
            $('#edit').val(true);
            $('#titlestr').val(data.title);
            $('#post').val(data.post);
            openForm("Edit Post #" + data.postID);
        });
    });

    $('.delete').click(function () {
        var id = $(this).data("id");
        if (confirm("Are you sure you want to remove this post?")) {
            $.post('/Forum/DeletePost', { 'id': id }, function (data) {
                if (data == "") {
                    $('#post_' + id).fadeOut('fast', function () { $('#post_' + id).remove(); });
                } else {
                    showMessage(data);
                }
            }, "text");
        }
    });

    $('.flag').click(function () {
        var id = $(this).data("id");
        if (confirm("Are you sure you want to remove the Flag as spam?")) {
            $.post('/Forum/UnFlagPost', { 'id': id }, function (data) {
                if (data == "") {
                    $('#post_' + id).fadeOut('fast', function () { $('#post_' + id).remove(); });
                } else {
                    showMessage(data);
                }
            }, "text");
        }
    });

    $('.block').click(function () {
        var id = $(this).data("id");
        var reason = prompt("Enter the reason for blocking this IP", "Spam");
        if (reason != null) {
            $.post('/Forum/BlockIP', { 'id': id, 'reason': reason }, function (data) {
                if (data == "") {
                    $('#post_' + id).fadeOut('fast', function () { $('#post_' + id).remove(); });
                } else {
                    showMessage(data);
                }
            }, "text");
        }
    });
});