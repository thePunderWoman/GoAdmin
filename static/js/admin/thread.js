$(function () {
    $('div.controls a').click(function (event) {
        event.preventDefault();
    });

    $('.approve').click(function () {
        var linkobj = $(this)
        var id = $(linkobj).data("id");
        $.post('/Forum/Approve', { 'id': id }, function (data) {
            $(linkobj).html(data);
        }, 'text');
    });

    $('.edit').click(function () {
        var id = $(this).data("id");
        $.getJSON('/Forum/GetPost', { 'postid': id }, function (data) {
            $('#postID').val(data.postID);
            $('#parentid').val(0);
            $('#edit').val(true);
            $('#titlestr').val(data.title);
            $('#post').val(data.post);
            if (data.sticky) {
                $('#sticky').attr('checked', true);
            }
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
        if (confirm("Are you sure you want to mark this post as spam?")) {
            $.post('/Forum/FlagPost', { 'id': id }, function (data) {
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

    $('#addreply').click(function (event) {
        event.preventDefault();
        $('#postID').val(0);
        $('#parentid').val(0);
        $('#edit').val(false);
        $('#titlestr').val('Re: ' + $('#index_title').find('span').html());
        $('#post').val('');
        openForm("Reply to Thread");
    });
});

function openForm(title) {
    $('#postform').dialog({
        autoOpen: false,
        title: title,
        width: 520,
        modal: true,
        buttons: {
            "Submit": function () {
                var sticky = ($('#sticky').is(':checked')) ? true : false;
                console.log(sticky);
                $.getJSON('/Forum/SavePost', { 'threadid': $('#threadID').val(), 'postid': $('#postID').val(), 'edit': $('#edit').val(), 'parentid': $('#parentID').val(), 'title': $('#titlestr').val(), 'post': $('#post').val(), 'sticky': sticky }, function (data) {
                    addPost(data.postID, $('#postID').val(), $('#edit').val())
                    clearForm();
                    $('#postform').dialog("close");
                });
            },
            Cancel: function () {
                $(this).dialog("close");
            }
        },
        close: function () { }
    });
    $('#postform').dialog('open');
}

function clearForm() {
    $('#postID').val(0);
    $('#parentid').val(0);
    $('#edit').val(false);
    $('#titlestr').val('');
    $('#post').val('');
    $('#sticky').attr('checked', false);
}

function addPost(postid, targetid, deleteObj) {
    if (targetid == undefined) {
        targetid = 0;
    }
    if (deleteObj == undefined) {
        deleteObj = false;
    }
    $.getJSON('/Forum/GetPost', { 'postid': postid }, function (data) {
        var email = data.email;
        var name = (data.name != "") ? data.name : "Anonymous";
        var post = '<div class="forumpost" id="post_' + data.postID + '">' +
                    '<div class="datebox">' +
                    '<p><strong>Posted On</strong><br />' + data.date + '</p>' +
                    '<p><strong>By</strong><br />' + ((email != "") ? "<a href='mailto:" + email + "'>" + name + "</a>" : name) +
                    ((data.company != "") ? '<br />' + data.company : '') + '</p>' +
                    '<p class="ipaddress"><strong>IP</strong> ' + data.IPAddress + '</p></div>' +
                    '<div class="postbox">' +
                    '<p class="title">' + data.title + '</p>' +
                    '<div class="controls">' +
                        '<a href="/Forum/Approve/' + data.postID + '" data-id="' + data.postID + '" class="approve">' + ((data.approved) ? "Unapprove" : "Approve") + '</a> |' +
                        '<a href="/Forum/EditPost/' + data.postID + '" data-id="' + data.postID + '" class="edit">Edit</a> | ' +
                        '<a href="/Forum/FlagPost/' + data.postID + '" data-id="' + data.postID + '" class="flag">Spam</a> | ' +
                        '<a href="/Forum/BlockIP/' + data.postID + '" data-id="' + data.postID + '" class="block">Block IP</a> | ' +
                        '<a href="/Forum/DeletePost/' + data.postID + '" data-id="' + data.postID + '" class="delete">Delete</a></div>' +
                    '<p>' + data.post + '</p></div><div class="clear"></div></div>';
        if (data.sticky) {
            $('#posts').children(':first').before(post);
            if (targetid != 0 && deleteObj) $($('#post_' + targetid)).remove();
        } else if (targetid == 0) {
            $('#posts').append(post);
        } else {
            var target = $('#post_' + targetid);
            $(target).after(post);
            if (deleteObj) $(target).remove();
        }
    });
}