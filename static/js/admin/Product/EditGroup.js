var groupTable, showForm, clearForm;

sortParts = function (parts) {
    parts.sort(function (a, b) {
        return (a.sort > b.sort) ? 1 : -1;
    });
    return parts;
}

buildPartList = function (parts) {
    parts = sortParts(parts);
    var partmsg = "";
    $(parts).each(function (i, part) {
        partmsg += '<li id="parts_' + part.id + '"><a target="_blank" href="/Product/Edit?partID=' + part.partID + '">' + part.partID + '</a><a class="removePart" href="/Product/RemovePartFromGroup/' + part.id + '">&times;</a>';
    });
    return partmsg;
}

updateGroupSort = (function () {
    var x = $('#groupPartList').sortable("serialize");
    $.post("/Product/updateGroupSort?" + x);
});

showForm = (function (groupID, name) {
    $('#groupID').val(groupID);
    $('#name').val(name);
    $('.form_left').slideDown();
});

clearForm = (function () {
    $('#groupID').val(0);
    $('#name').val('');
    $('.form_left').slideUp();
});

$(function () {
    groupTable = $('table').dataTable({ "bJQueryUI": true });

    $(document).on('click', '#addGroup', function (e) {
        e.preventDefault();
        showForm(0, '');
    });

    $(document).on('click','.edit', function (e) {
        e.preventDefault();
        var groupID = $(this).data('id');
        var clicked_link = $(this);
        $.getJSON('/Product/GetGroup', { 'groupID': groupID }, function (response) {
            if (response.error == null) {
                groupTable.fnDeleteRow($(clicked_link).parent().parent().get()[0]);
                showForm(response.id, response.name);
            } else {
                showMessage(response.error);
            }
        });
    });

    $(document).on('click', '.parts', function (e) {
        e.preventDefault();
        var groupID = $(this).data('id');
        $.getJSON('/Product/GetGroup', { 'groupID': groupID }, function (response) {
            console.log(response);
            $("#config-dialog").empty();
            var partmsg = '<p>Drag and drop to change order</p><ul id="groupPartList">';
            partmsg += buildPartList(response.Parts);
            partmsg += '</ul>';
            if (response.Parts.length == 0) {
                partmsg += '<p id="noparts">No Parts Associated</p>';
            }
            partmsg += '<label for="addPart">Add Part<br /><input type="text" id="addPart" data-id="' + groupID + '" placeholder="Enter a part number" /></label>';
            partmsg += '<button id="submitPart">Add</button>';
            $("#config-dialog").append(partmsg);
            $("#config-dialog").dialog({
                modal: true,
                title: "Group Parts",
                width: 'auto',
                height: 'auto',
                buttons: {
                    "Done": function () {
                        $('#groupPartList').sortable("destroy");
                        $(this).dialog("close");
                    }
                }
            });
            $('#groupPartList').sortable({ axis: "y", update: function (event, ui) { updateGroupSort(event, ui) } }).disableSelection();
        });
    });
    

    $(document).on('click', '.remove', function (e) {
        e.preventDefault();
        clearForm();
        var clicked_link = $(this);
        var groupID = $(this).data('id');
        if (groupID > 0 && confirm('Are you sure you want to delete this group?')) {
            $.get('/Product/DeleteGroup', { 'groupID': groupID }, function (response) {
                if (response == "") {
                    groupTable.fnDeleteRow($(clicked_link).parent().parent().get()[0]);
                    showMessage("Group removed.");

                } else {
                    showMessage(response);
                }
            });
        } else if (contentID <= 0) {
            showMessage("Group ID not valid.");
        }
    });

    $(document).on('click','#btnReset', function () {
        var groupID = $('#groupID').val();
        if (groupID > 0) {
            $.getJSON('/Product/GetGroup', { 'groupID': groupID }, function (response) {
                var addId = groupTable.fnAddData([
                                response.name,
                                response.Parts.length,
                                '<a href="#" class="edit" data-id="' + response.id + '" title="Edit Group">Edit</a> | <a href="#" class="parts" data-id="' + response.id + '" title="Edit Parts">Edit Parts</a> | <a href="#" class="remove" data-id="' + response.id + '" title="Remove Group">Remove</a>'
                ]);
                var theCell = groupTable.fnSettings().aoData[addId[0]].nTr.cells[2];
                theCell.className = "center"
            })
        }
        clearForm();
    });

    $(document).on('click','#btnSave', function () {
        var name = $('#name').val().trim();
        var partID = $('#partID').val();
        if (partID > 0 && name.length > 0) {
            var groupID = $('#groupID').val();
            $.getJSON("/Product/SaveGroup", { 'partID': partID, 'name': name, 'groupID': groupID }, function (response) {
                if (response.error == null) {
                    var addId = groupTable.fnAddData([
                                    response.name,
                                    response.Parts.length,
                                    '<a href="#" class="edit" data-id="' + response.id + '" title="Edit Group">Edit</a> | <a href="#" class="parts" data-id="' + response.id + '" title="Edit Parts">Edit Parts</a> | <a href="#" class="remove" data-id="' + response.id + '" title="Remove Group">Remove</a>'
                    ]);
                    var theCell = groupTable.fnSettings().aoData[addId[0]].nTr.cells[2];
                    theCell.className = "center"
                    showMessage("Group Saved.");
                    clearForm();
                } else {
                    showMessage(response.error);
                }
            });
        } else {
            if (partID <= 0) {
                showMessage("Error getting part number.");
            } else if (name.length == 0) {
                showMessage("Name cannot be blank.");
            } else {
                showMessage("Error encountered.");
            }
        }
        return false;
    });

    $(document).on('click', '#submitPart', function (e) {
        e.preventDefault();
        var bobj = $(this);
        var partID = $('#addPart').val().trim();
        if (partID != "") {
            var groupID = $('#addPart').data('id');
            $.post('/Product/AddGroupPart', { groupID: groupID, partID: partID }, function (data) {
                $('#groupPartList').sortable("destroy")
                $('#noparts').remove();
                $('#groupPartList').empty();
                var partmsg = buildPartList(data.Parts);
                $('#groupPartList').append(partmsg);
                $('#groupPartList').sortable({ axis: "y" ,update: function (event, ui) { updateGroupSort(event, ui) } }).disableSelection();
                $('#addPart').attr('value', '');
            },"json");
        } else {
            $('#addPart').attr('value', '');
            showMessage("You must enter a part ID.");
        }
    });

    $(document).on('click', '.removePart', function (e) {
        e.preventDefault();
        var href = $(this).attr('href');
        var liobj = $(this).parent();
        if (confirm('Are you sure you want to remove this part from this group?')) {
            $.post(href, function (data) {
                if (data) {
                    $(liobj).fadeOut('400', function () {
                        $(liobj).remove();
                        if ($('#groupPartList li').length == 0) {
                            $('#groupPartList').after('<p id="noparts">No Parts Associated</p>');
                        }
                    });
                }
            }, "json");
        }
    });
});