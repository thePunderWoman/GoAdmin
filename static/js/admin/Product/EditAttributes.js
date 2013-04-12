var showAttributeForm, clearAttributeForm, addRow, updateRow, updateAttributeSort;

showAttributeForm = (function (field, value, attrID) {
    $('#field').val(field);
    $('#value').val(value);
    $('#attrID').val(attrID);
    $('form.form_left').slideDown();
});

clearAttributeForm = (function () {
    $('#field').val('');
    $('#value').val('');
    $('#attrID').val(0);
    $('form.form_left').slideUp();
});

addRow = (function (data) {
    var newrow = '<tr class="sortable" id="attr_' + data.pAttrID + '">';
    newrow += '<td class="handles">' + data.field + '</td><td class="handles">' + data.value + '</td><td><a href="javascript:void(0)" class="edit" id="' + data.pAttrID + '" title="Edit ' + data.field + '">Edit</a> | <a href="javascript:void(0)" class="remove" id="' + data.pAttrID + '" title="Remove ' + data.field + '">Remove</a></td>';
    newrow += '</tr>';
    $('table tbody').append(newrow);
});

updateRow = (function (idstr, data) {
    var urow = $('#attr_' + idstr);
    var urowcontent = '<td class="handles">' + data.field + '</td><td class="handles">' + data.value + '</td><td><a href="javascript:void(0)" class="edit" id="' + data.pAttrID + '" title="Edit ' + data.field + '">Edit</a> | <a href="javascript:void(0)" class="remove" id="' + data.pAttrID + '" title="Remove ' + data.field + '">Remove</a></td>';
    $(urow).empty();
    $(urow).append(urowcontent);
})

updateAttributeSort = (function () {
    var x = $('table tbody').sortable("serialize");
    $.post("/Product/updateAttributeSort?" + x);
});

$(function () {
    var fixHelper = function (e, ui) {
        ui.children().each(function () {
            $(this).width($(this).width());
        });
        return ui;
    };
    $('table tbody').sortable({ helper: fixHelper, handle: 'td.handles', update: function (event, ui) { updateAttributeSort(event, ui) } }).disableSelection();

    $(document).on('click', '#addField', function () {
        var field = prompt('Please enter the new attribute field name.');
        field = $.trim(field);
        $('#field').append('<option>' + field + '</option>');
        $('#field').val(field);
    });

    $(document).on('click', '#addAttribute', function () {
        showAttributeForm('', '', 0);
    });

    $(document).on('click', '.edit', function () {
        var attrID = $(this).attr('id');
        var field = $(this).parent().prev().prev().text();
        var value = $(this).parent().prev().text();
        showAttributeForm(field, value, attrID);
    });

    $(document).on('click', '.remove', function () {
        clearAttributeForm();
        var clicked_link = $(this);
        var attrID = $(this).attr('id');
        if (attrID > 0 && confirm('Are you sure you want to remove this attribute?')) {
            $.get('/Product/DeleteAttribute', { 'attrID': attrID }, function (response) {
                $('#attr_' + attrID).fadeOut('fast', function () { $('#attr_' + attrID).remove(); });
                showMessage(response);
            });
        } else if (attrID.length <= 0 || attrID <= 0) {
            showMessage("Attribute ID invalid.");
        }
    });

    $(document).on('click', '#btnSave', function () {
        var field, value, partID;
        field = $('#field').val();
        value = $('#value').val();
        partID = $('#partID').val();
        if (partID > 0 && field.length > 0 && value.length > 0) {
            var attrID = $('#attrID').val();
            $.getJSON('/Product/SaveAttribute', { 'attrID': attrID, 'partID': partID, 'field': field, 'value': value }, function (response) {
                if (response.error == null) {
                    if (attrID == 0) {
                        addRow(response);
                        $('table tbody').sortable("destroy");
                        $('table tbody').sortable({ helper: fixHelper, handle: 'td.handles', update: function (event, ui) { updateAttributeSort(event, ui) } }).disableSelection();
                    } else {
                        updateRow(attrID, response);
                    }
                    clearAttributeForm();
                    showMessage('Attribute Saved.');
                } else {
                    showMessage(response.error);
                }
            });
        } else {
            if (partID <= 0) { showMessage("Can't find Part ID."); }
            if (field.length == 0) { showMessage("You must select a field."); }
            if (value.length == 0) { showMessage("You must enter a value."); }
        }
        return false;
    });

});

