var selectFile, imageSort;

$(function () {
    $("#tabs").tabs();

    $('#startDate,#endDate').datetimepicker({
        ampm: true
    });

    $('#addImage').on('click', function (e) {
        e.preventDefault();
        chooseFile();
    });

    $('#addData').on('click', function (e) {
        e.preventDefault();
        var pageID = $('#pageID').val();
        var key = $('#dataKey').val().trim();
        var value = $('#dataValue').val().trim();
        if (key == '' || value == '') {
            showMessage('You must have a key and a value to add a custom data attribute.');
        } else {
            $.getJSON('/LandingPages/AddData', { pageID: pageID, key: key, value: value }, function (data) {
                $('#dataList').empty();
                var datalist = "";
                $(data).each(function (i, obj) {
                    datalist += '<li>' + obj.dataKey + ': ' + obj.dataValue + ' <a href="/LandingPages/RemoveData/' + obj.id + '" class="removeData">&times;</a></li>';
                });
                $('#dataList').append(datalist);
            });
            $('#dataKey').val('');
            $('#dataValue').val('');
        }
    });

    imageSort();

    $(document).on('click', 'a.removeimage', function (e) {
        e.preventDefault();
        var href = $(this).attr('href');
        var liobj = $(this).parent();
        $.getJSON(href, function (data) {
            $(liobj).fadeOut(400, function () {
                $(this).remove();
            });
        });
    });

    $(document).on('click', 'a.removeData', function (e) {
        e.preventDefault();
        var href = $(this).attr('href');
        var liobj = $(this).parent();
        $.getJSON(href, function (data) {
            $(liobj).fadeOut(400, function () {
                $(this).remove();
            });
        });
    });
    

    CKEDITOR.replace('page_content', {
        filebrowserImageUploadUrl: '/File/CKUpload',
        filebrowserImageBrowseUrl: '/File/CKIndex',
        filebrowserImageWindowWidth: '640',
        filebrowserImageWindowHeight: '480'
    });

});

selectFile = function(url) {
    var pageID = $('#pageID').val();
    $.getJSON('/LandingPages/AddImage', { pageID: pageID, image: url }, function (data) {
        $('#pageImages').empty();
        var lis = "";
        $(data).each(function (i, obj) {
            lis += '<li id="img_' + obj.id + '"><img src="' + obj.url + '" alt="page image" /><a href="/LandingPages/RemoveImage/' + obj.id + '" class="removeimage">&times;</a><span class="clear"></span></li>';
        });
        $('#pageImages').append(lis);
    });
    $("#file-dialog").dialog("close");
    $("#file-dialog").empty();
}

imageSort = function () {
    $("#pageImages").sortable("destroy")
    $("#pageImages").sortable({
        handle: "img",
        axis: "y",
        cursor: "move",
        update: function (event, ui) {
            var pageID = $('#pageID').val();
            var sortstr = $("#pageImages").sortable("serialize", { key: "img" });
            $.post('/LandingPages/UpdateSort?' + sortstr);
        }
    });
    $("#pageImages").disableSelection();
}