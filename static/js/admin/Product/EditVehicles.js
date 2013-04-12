var relatedTable = "";
var allTable = "";
var attrTable = "";
var fixHelper = function (e, ui) {
    ui.children().each(function () {
        $(this).width($(this).width());
    });
    return ui;
};
$(function () {
    relatedTable = $('#relatedVehicles').dataTable({ "bJQueryUI": true, "aaSorting": [[5, "asc"]] });
    $('#tableContainer').fadeIn();
    $('#addvehicle').hide();

    /*$.getJSON('/Product/getAllVehicles', function (data) {
    $.each(data, function (i, obj) {
    $('#allVehiclesBody').append('<tr id="allVehicle:' + obj.vehicleID + '"><td>' + obj.vehicleID + '</td><td>' + obj.year + '</td><td>' + obj.make + '</td><td>' + obj.model + '</td><td>' + obj.style + '</td><td><a href="javascript:void(0)" class="add" id="' + obj.vehicleID + '">Add</a></td></tr>');
    });
    allTable = $('#allVehicles').dataTable({ "bJQueryUI": true });
    $('#loading_area').fadeOut(function () { $('#allVehicles').fadeIn(); });
    });*/

    $('#addField').live('click', function () {
        var field = prompt('Please enter the new attribute field name.');
        $('#field').append('<option>' + field + '</option>');
        $('#field').val(field);
    });

    $(document).on('click', '#addvehicle', function () {
        var yearID = $('#select_year').val();
        var makeID = $('#select_make').val();
        var modelID = $('#select_model').val();
        var styleID = $('#select_style').val();
        var partID = $('#partID').val();
        if (yearID > 0 && makeID > 0 && modelID > 0 && styleID > 0 && partID > 0) {
            $.getJSON('/Product/GetAllPartOptions', { 'partID': partID }, function (data) {
                var partstr = "";
                if (data.length == 1) {
                    partstr = " #" + data[0];
                } else if (data.length == 2) {
                    partstr = "s #" + data[0] + " and #" + data[1];
                } else if (data.length > 2) {
                    $(data).each(function (i, obj) {
                        if (i == 0) {
                            partstr += "s #" + obj;
                        } else if (i != (data.length - 1)) {
                            partstr += ", #" + obj;
                        } else {
                            partstr += " and #" + obj;
                        }
                    });
                }

                if (confirm('Are you sure you want to add ' + $("#select_year option[value='" + yearID + "']").text() + ' ' + $("#select_make option[value='" + makeID + "']").text() + ' ' + $("#select_model option[value='" + modelID + "']").text() + ' ' + $("#select_style option[value='" + styleID + "']").text() + ' vehicle to part' + partstr + '?')) {
                    var editing = $('#edit').val();

                    $.getJSON('/Product/AddVehicleByFilter', { 'partID': partID, 'yearID': yearID, 'makeID': makeID, 'modelID': modelID, 'styleID': styleID }, function (response) {
                        if (response.error == null && response.vehicleID > 0) {
                            //allTable.fnDeleteRow($('#allVehicle\\:' + vehicleID).get()[0]);
                            var addId = relatedTable.fnAddData([
                                    response.vehicleID,
                                    response.year,
                                    response.make,
                                    response.model,
                                    response.style,
                                    response.make + " " + response.model + " " + response.style,
                                    '<a href="javascript:void(0)" id="' + response.vehicleID + '" class="edit">Edit</a> | <a href="javascript:void(0)" id="' + response.vehicleID + '" class="carryover">Carry Over</a> | <a href="javascript:void(0)" class="remove" id="' + response.vehicleID + '">Remove</a>'
                                ]);

                            var theNode = relatedTable.fnSettings().aoData[addId[0]].nTr;
                            theNode.setAttribute('id', 'relatedVehicle_' + vehicleID);

                            showMessage('Vehicle added.');
                            editVehicle(vehicleID)
                        } else {
                            showMessage(response.error);
                        }
                    });
                }
            });
        }
    });

    $('.edit').live('click', function () {
        var vehicleID = $(this).attr('id');
        editVehicle(vehicleID)
    });

    $('.carryover').live('click', function () {
        var vehicleID = $(this).attr('id');
        var partID = $('#partID').val();
        $.getJSON('/Product/GetCarryOverData', { 'vehicleID': vehicleID, 'partID': partID }, function (data) {
            var partstr = "";
            if (data.partids.length == 1) {
                partstr = " #" + data.partids[0];
            } else if (data.partids.length == 2) {
                partstr = "s #" + data.partids[0] + " and #" + data.partids[1];
            } else if (data.partids.length > 2) {
                $(data.partids).each(function (i, obj) {
                    if (i == 0) {
                        partstr += "s #" + obj;
                    } else if (i != (data.partids.length - 1)) {
                        partstr += ", #" + obj;
                    } else {
                        partstr += " and #" + obj;
                    }
                });
            }
            if (confirm("This will carry part" + partstr + " over to the following vehicle:\n\n" + (data.year + 1) + " " + data.make + " " + data.model + " " + data.style + "\n\n Are you sure you want to do this?")) {
                $.post('/Product/CarryOverPart', { 'vehicleID': vehicleID, 'partID': partID }, function (response) {
                    if (response.error == null && response.vehicleID > 0) {
                        var addId = relatedTable.fnAddData([
                            response.vehicleID,
                            response.year,
                            response.make,
                            response.model,
                            response.style,
                            response.make + " " + response.model + " " + response.style,
                            '<a href="javascript:void(0)" id="' + response.vehicleID + '" class="edit">Edit</a> | <a href="javascript:void(0)" id="' + response.vehicleID + '" class="carryover">Carry Over</a> | <a href="javascript:void(0)" class="remove" id="' + response.vehicleID + '">Remove</a>'
                        ]);

                        var theNode = relatedTable.fnSettings().aoData[addId[0]].nTr;
                        theNode.setAttribute('id', 'relatedVehicle_' + vehicleID);

                        showMessage('Part Carried Over.');
                        //editVehicle(vehicleID)
                    } else {
                        showMessage(response.error);
                    }
                }, "json");
            }
        });
        //editVehicle(vehicleID)
    });

    $('.remove').live('click', function () {
        var vehicleID = $(this).attr('id');
        var partID = $('#partID').val();
        var clicked_link = $(this);
        if (vehicleID > 0 && partID > 0 && confirm('Are you sure you want to remove the relationship to Vehicle #' + vehicleID + '?')) {
            $.get('/Product/DeleteVehicle', { 'vehicleID': vehicleID, 'partID': partID }, function (response) {
                if (response.length == 0) {
                    relatedTable.fnDeleteRow($(clicked_link).parent().parent().get()[0]);
                    showMessage("Vehicle removed from part.");
                } else {
                    showMessage(response);
                }
            });
        }
    });

    $('.editattr').live('click', function (event) {
        event.preventDefault();
        var attrID = $(this).attr('id').split('_')[1];
        $.getJSON('/Product/GetVehiclePartAttribute', { 'vpAttrID': attrID }, function (response) {
            $('#field').val(response.field);
            $('#value').val(response.value);
            $("form.form_left").dialog({
                autoOpen: false,
                height: 300,
                width: 620,
                title: "Edit Attribute",
                modal: true,
                buttons: {
                    "Save": function () {
                        var vpAttrID = $('#vpAttrID').val();
                        var field = $('#field').val();
                        var val = $('#value').val();
                        var bValid = true;

                        if ($.trim(field) == "") bValid = false;
                        if ($.trim(val) == "") bValid = false;

                        if (bValid) {
                            $.post('/Product/UpdateVehiclePartAttribute', { vpAttrID: vpAttrID, field: field, value: val }, function (data) {
                                // remove and re-add row to table
                                var removeme = $('#attribute_' + data.vpAttrID);
                                $(removeme).before(addRow(data.vpAttrID, data.field, data.value));
                                $(removeme).remove();
                                $('#attributeTable tbody').sortable("destroy");
                                $('#attributeTable tbody').sortable({ helper: fixHelper, handle: 'td.handles', update: function (event, ui) { updateAttributeSort(event, ui) } }).disableSelection();
                            }, "json");
                            $(this).dialog("close");
                        }
                    },
                    Cancel: function () {
                        $(this).dialog("close");
                    }
                },
                close: function () { }
            });
            $("form.form_left").dialog("open");
        });
        $('#vpAttrID').val(attrID);
    });

    $('.removeattr').live('click', function () {
        var idstr = $(this).attr('id').split('_')[1];
        if (confirm('Are you sure you want to remove the Information relating to this Part and Vehicle?')) {
            $.get('/Product/DeleteVehiclePartAttribute', { 'vpAttrID': idstr }, function (response) {
                if (response.length == 0) {
                    $('#attribute_' + idstr).remove();
                    showMessage("Vehicle Part Attribute removed.");
                } else {
                    showMessage(response);
                }
            });
        }
    });

    $('#addAttribute').live('click', function (event) {
        event.preventDefault();
        var attrID = $(this).attr('id').split('_')[1];
        $('#field').val('');
        $('#value').val('');
        $("form.form_left").dialog({
            autoOpen: false,
            height: 300,
            width: 620,
            title: "Add Attribute",
            modal: true,
            buttons: {
                "Save": function () {
                    var vPartID = $('#vPartID').val();
                    var field = $('#field').val();
                    var val = $('#value').val();
                    var bValid = true;

                    if ($.trim(field) == "") bValid = false;
                    if ($.trim(val) == "") bValid = false;

                    if (bValid) {
                        $.post('/Product/AddVehiclePartAttribute', { vPartID: vPartID, field: field, value: val }, function (data) {
                            // remove and re-add row to table
                            $('#attributeTable tbody').append(addRow(data.vpAttrID, data.field, data.value));
                            $('#attributeTable tbody').sortable("destroy");
                            $('#attributeTable tbody').sortable({ helper: fixHelper, handle: 'td.handles', update: function (event, ui) { updateAttributeSort(event, ui) } }).disableSelection();
                        }, "json");
                        $(this).dialog("close");
                    }
                },
                Cancel: function () {
                    $(this).dialog("close");
                }
            },
            close: function () { }
        });
        $("form.form_left").dialog("open");
        $('#vpAttrID').val(attrID);
    });

    $('#btnReset').live('click', function () {
        $('#field').val('');
        $('#attributes').slideUp();
        $('#attributeTable tbody').sortable("destroy");
    });

    // Handle the change of year
    $(document).on('change', '#select_year', function () {
        // Initiate make, model, and style to default state
        $('#lookupError').hide();
        var yearID = $('#select_year').val();
        $('#select_make').html('<option value="">- Select Make -</option>');
        $('#select_model').html('<option value="">- Select Model -</option>');
        $('#select_model').attr('disabled', 'disabled');
        $('#select_style').html('<option value="">- Select Style -</option>');
        $('#select_style').attr('disabled', 'disabled');
        $('#addvehicle').hide();

        if (yearID > 0) {
            $.getJSON('/Vehicles/GetMakes', { 'yearID': yearID }, loadMake);
        }
    });

    // Handle the change of make
    $(document).on('change', '#select_make', function () {
        // Initiate make, model, and style to default state
        $('#lookupError').hide();
        var yearID = $('#select_year').val();
        var makeID = $('#select_make').val();
        $('#select_model').html('<option value="">- Select Model -</option>');
        $('#select_style').html('<option value="">- Select Style -</option>');
        $('#select_style').attr('disabled', 'disabled');
        $('#addvehicle').hide();
        if (yearID > 0 && makeID > 0) {
            $.getJSON('/Product/GetModels', { 'yearID': yearID, 'makeID': makeID }, loadModel);
        }
    });

    // Handle the change of model
    $(document).on('change', '#select_model', function () {
        $('#lookupError').hide();
        var yearID = $('#select_year').val();
        var makeID = $('#select_make').val();
        var modelID = $('#select_model').val();
        $('#addvehicle').hide();
        $('#select_style').html('<option value="">- Select Style -</option>');
        if (yearID > 0 && makeID > 0 && modelID > 0) {
            $.getJSON('/Product/GetStyles', { 'yearID': yearID, 'makeID': makeID, 'modelID': modelID }, loadStyle);
        }
    });

    $(document).on('change', '#select_style', function () {
        var yearID = $('#select_year').val();
        var makeID = $('#select_make').val();
        var modelID = $('#select_model').val();
        var styleID = $('#select_style').val();
        if (yearID > 0 && makeID > 0 && modelID > 0 && styleID > 0) {
            $('#addvehicle').fadeIn();
        } else {
            $('#addvehicle').hide();
        }
    });

});


function editVehicle(vehicleID) {
    var partID = $('#partID').val();
    var tds = $('#relatedVehicle_' + vehicleID).find('td');
    var header = 'Editing ' + $(tds[1]).text() + ' ' + $(tds[2]).text() + ' ' + $(tds[3]).text() + ' ' + $(tds[4]).text();
    if (vehicleID > 0 && partID > 0) {
        $.getJSON('/Product/GetVehiclePart', { 'vehicleID': vehicleID, 'partID': partID }, function (response) {
            if (response.error == null) {
                $('#vehicleID').val(response.vehicleID);
                $('#vPartID').val(response.vPartID);
                $('#attributeTable').find('tbody').empty();
                $.each(response.attributes, function (i, obj) {
                    $('#attributeTable tbody').append(addRow(obj.vpAttrID, obj.field, obj.value));
                });
                $('#attributes').find('h4').text(header);
                $('#attributes').slideDown();
                $('#attributeTable tbody').sortable("destroy");
                $('#attributeTable tbody').sortable({ helper: fixHelper, handle: 'td.handles', update: function (event, ui) { updateAttributeSort(event, ui) } }).disableSelection();
            } else {
                showMessage(response.error);
            }
        });
    }
}

function addRow(idstr, field, value) {
    var row = '<tr id="attribute_' + idstr + '"><td class="handles">' + field + '</td><td class="handles">' + value + '</td><td><a href="javascript:void(0)" id="editattr_' + idstr + '" class="editattr">Edit</a> | <a href="javascript:void(0)" id="removeattr_' + idstr + '" class="removeattr">Remove</a></td></tr>'
    return row;
}

function loadMake(makes) {
    $(makes).each(function (i, make) {
        $('#select_make').append('<option value="' + make.makeID + '">' + make.make1 + '</option>');
    });
    $('#select_make').removeAttr('disabled');
}

function loadModel(models) {
    $(models).each(function (i, model) {
        $('#select_model').append('<option value="' + model.modelID + '">' + model.model1 + '</option>');
    });
    $('#select_model').removeAttr('disabled');
}

function loadStyle(styles) {
    $(styles).each(function (i, style) {
        $('#select_style').append('<option value="' + style.styleID + '">' + style.style1 + '</option>');
    });
    $('#select_style').removeAttr('disabled');
}

function updateAttributeSort() {
    var x = $('table tbody').sortable("serialize");
    $.post("/Product/updateVehicleAttributeSort?" + x);
}