var contentTable = "";
$(function () {
    var contentTable = $('table').dataTable({ "bJQueryUI": true });

    $('.action').live('change', function () {
        var contentID = $(this).attr('id');
        var catID = $('#categoryID').val();
        var action = $(this).val();

        switch (action) {
            case 'edit':
                $.getJSON('/Categories/GetContent', { 'contentID': contentID }, function (response) {
                    $('#contentID').val(response.contentID);
                    $('#content').val(response.content);
                    $('#contentType').val(response.content_type_id).trigger('change');
                    $('form.form_left').slideDown();
                });
                break;

            case 'delete':
                // Delete this category
                if (confirm("Are you sure you want to remove this content?\r\nThis cannot be undone!")) {
                    $.getJSON('/Categories/DeleteContent', { 'catID': catID, 'contentID': contentID }, function (response) {
                        if ($.trim(response).length == 0) {
                            contentTable.fnDeleteRow($('#contentRow\\:' + contentID).get()[0]);
                            showMessage('Content has been removed from this Category.');
                        } else {
                            showMessage('' + response);
                        }
                    });
                }
                break;

            default:
                break;
        }
        $(this).val(0);
    });

    $('#addContent').live('click', function (event) {
        event.preventDefault();
        $('form.form_left').slideDown();
    });

    $('#contentType').live('change', function (event) {
        var cTypeID = $(this).val();
        if (cTypeID != '') {
            $.getJSON('/Misc/GetContentType', { 'cTypeID': cTypeID }, function (data) {
                var content = $('#content').val();
                if (data.allowHTML) {
                    CKEDITOR.replace('content', {
                        filebrowserImageUploadUrl: '/File/CKUpload',
                        filebrowserImageBrowseUrl: '/File/CKIndex'
                    });
                } else {
                    if (CKEDITOR.instances.content != undefined) CKEDITOR.instances.content.destroy();
                }
            });
        }
    });

    $('#btnSave').click(function (event) {
        event.preventDefault();
        var contentID = $('#contentID').val();
        var catID = $('#categoryID').val();
        var content = (CKEDITOR.instances.content == undefined) ? $('#content').val() : CKEDITOR.instances.content.getData();
        var typeID = $('#contentType').val();
        $.post('/Categories/SaveContent', { 'catID': catID, 'contentID': contentID, 'typeID': typeID, 'content': content }, function (data) {
            if (contentID != 0) {
                contentTable.fnDeleteRow($('#contentRow\\:' + data.contentID).get()[0]);
            }

            var addId = contentTable.fnAddData([
                                        data.content_type,
                                        '',
                                        '<select class="action" id="' + data.contentID + '"><option value="0">- Select Option -</option><option value="edit">Edit</option><option value="delete">Delete</option></select>'
                                    ]);

            var theNode = contentTable.fnSettings().aoData[addId[0]].nTr;
            theNode.setAttribute('id', 'contentRow\:' + data.contentID);
            var theCell = contentTable.fnSettings().aoData[addId[0]].nTr.cells[1];
            var contentstr = (content.length <= 90) ? content : content.substr(0, 90) + "...";
            $(theCell).text(contentstr);
            var theCell = contentTable.fnSettings().aoData[addId[0]].nTr.cells[2];
            theCell.className = "center"

            showMessage('Content saved.');

        }, "json");
        $('#btnReset').trigger('click');
    });

    $('#btnReset').click(function (event) {
        $('#contentType').val('');
        $('#contentID').val(0);
        $('#content').val('');
        if (CKEDITOR.instances.content != undefined) CKEDITOR.instances.content.destroy();
        $('form.form_left').slideUp();
    });

});