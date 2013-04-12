var contentTable, showContentForm, clearContentForm;
showContentForm = (function (contentID, contentType, content) {
    $('#contentID').val(contentID);
    $('#contentType').val(contentType).trigger('change');
    $('#content').val(content);
    $('.form_left').slideDown();
});

clearContentForm = (function () {
    $('#contentType').val('');
    $('#contentID').val(0);
    $('#content').val('');
    if (CKEDITOR.instances.content != undefined) CKEDITOR.instances.content.destroy();
    $('.form_left').slideUp();
});

$(function () {
    contentTable = $('table').dataTable({ "bJQueryUI": true });

    $(document).on('click','#addContent', function () {
        showContentForm(0, '', '');
        $('#editing').val(0);
    });

    $(document).on('click','.edit', function () {
        var contentID = $(this).attr('id');
        var clicked_link = $(this);
        $.getJSON('/Product/GetFullContent', { 'contentID': contentID }, function (response) {
            if (response.error == null) {
                contentTable.fnDeleteRow($(clicked_link).parent().parent().get()[0]);
                showContentForm(response.contentID, response.content_type_id, response.content);
            } else {
                showMessage(response.error);
            }
        });
    });

    $(document).on('change','#contentType', function (event) {
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

    $(document).on('click','.remove', function () {
        clearContentForm();
        var clicked_link = $(this);
        var partID = $('#partID').val();
        var contentID = $(this).attr('id');
        if (contentID > 0 && confirm('Are you sure you want to remove this content?')) {
            $.get('/Product/DeleteContent', { 'partID': partID, 'contentID': contentID }, function (response) {
                if (response == "") {
                    contentTable.fnDeleteRow($(clicked_link).parent().parent().get()[0]);
                    showMessage("Content removed.");

                } else {
                    showMessage(response);
                }
            });
        } else if (contentID <= 0) {
            showMessage("Content ID not valid.");
        }
    });

    $(document).on('click','#btnReset', function () {
        var contentID = $('#contentID').val();
        if (contentID > 0) {
            var content_type_id = $('#contentType').val();
            var content_type = $("#contentType option[value='" + content_type_id + "']").text();
            var content = $('#content').val();
            var addId = contentTable.fnAddData([
                            content_type,
                            content,
                            '<a href="javascript:void(0)" class="edit" id="' + contentID + '" title="Edit Content">Edit</a> | ' + '<a href="javascript:void(0)" class="remove" id="' + contentID + '" title="Remove Content">Remove</a>'
                        ]);
            var theNode = contentTable.fnSettings().aoData[addId[0]].nTr;
            theNode.setAttribute('id', 'contentRow\:' + contentID);
            var theCell = contentTable.fnSettings().aoData[addId[0]].nTr.cells[1];
            $(theCell).text(content);
            var theCell = contentTable.fnSettings().aoData[addId[0]].nTr.cells[2];
            theCell.className = "center"
            showMessage("Content added.");
        }
        clearContentForm();
    });

    $(document).on('click','#btnSave', function () {
        var contentType = $('#contentType').val();
        var content = (CKEDITOR.instances.content == undefined) ? $('#content').val() : CKEDITOR.instances.content.getData();
        var partID = $('#partID').val();
        if (partID > 0 && content.length > 0 && contentType > 0) {
            var contentID = $('#contentID').val();
            $.getJSON("/Product/SaveContent", { 'contentID': contentID, 'partID': partID, 'content': content, 'contentType': contentType }, function (response) {
                if (response.error == null) {
                    var addId = contentTable.fnAddData([
                                    response.content_type,
                                    response.content,
                                    '<a href="javascript:void(0)" class="edit" id="' + response.contentID + '" title="Edit Content">Edit</a> | ' + '<a href="javascript:void(0)" class="remove" id="' + response.contentID + '" title="Remove Content">Remove</a>'
                                ]);
                    var theNode = contentTable.fnSettings().aoData[addId[0]].nTr;
                    theNode.setAttribute('id', 'contentRow\:' + response.contentID);
                    var theCell = contentTable.fnSettings().aoData[addId[0]].nTr.cells[1];
                    $(theCell).text(content);
                    var theCell = contentTable.fnSettings().aoData[addId[0]].nTr.cells[2];
                    theCell.className = "center"
                    showMessage("Content Saved.");
                    clearContentForm();
                } else {
                    showMessage(response.error);
                }
            });
        } else {
            if (partID <= 0) {
                showMessage("Error getting part number.");
            } else if (content.length == 0) {
                showMessage("Content cannot be blank.");
            } else if (contentType <= 0) {
                showMessage("You must select a content type.");
            } else {
                showMessage("Error encountered.");
            }
        }
        return false;
    });
});