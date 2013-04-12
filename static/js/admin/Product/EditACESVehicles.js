var getPartVehicles, generatePartConfigTable;
$(function () {
    $("#tabs").tabs();
    $('#find').hide();

    $('div.configs').show();

    $(document).on('click', '.mapPart', function (e) {
        e.preventDefault();
        var id = $(this).data('id');
        if (confirm('This will populate the ACES vehicles with best guesses from the non-aces vehicles mapped to this part. Are you sure?')) {
            $.post('/ACES/MapPart/' + id, function (data) {
                getPartVehicles();
                if (data.length > 0) {
                    var unmerged = "<p>The Following vehicles did not map:</p><ul>";
                    $(data).each(function (i, obj) {
                        unmerged += "<li>" + obj.Year.year1 + " " + obj.Make.make1 + " " + obj.Model.model1 + " " + obj.Style.style1 + "</li>";
                    });
                    unmerged += "</ul>";
                    $("#config-dialog").append(unmerged);
                    $("#config-dialog").dialog({
                        modal: true,
                        title: "Vehicles That Failed to Map",
                        width: 'auto',
                        height: 'auto',
                        buttons: {
                            "OK": function () {
                                $(this).dialog("close");
                                $("#config-dialog").empty();
                            }
                        }
                    });
                }
            }, "json");
        }
    });

    $(document).on('click', '.removeBV,.removeSubmodel', function (e) {
        e.preventDefault();
        var href = $(this).attr('href');
        var toolobj = $(this).parent();
        if (confirm('Are you sure you want to remove this vehicle from this part?')) {
            $.post(href, function (data) {
                if (data) {
                    $(toolobj).fadeOut('400', function () {
                        $(toolobj).remove();
                    });
                    getCurtDevVehicles();
                }
            }, "json");
        }
    });

    $(document).on('click', '.removeConfig', function (e) {
        e.preventDefault();
        var href = $(this).attr('href');
        var trobj = $(this).parent().parent();
        if (confirm('Are you sure you want to remove this vehicle from this part?')) {
            $.post(href, function (data) {
                if (data) {
                    $(trobj).fadeOut('400', function () {
                        $(trobj).remove();
                    });
                    getCurtDevVehicles();
                }
            }, "json");
        }
    });

    $(document).on('click', '.addToPart', function (e) {
        e.preventDefault();
        var aobj = $(this);
        var href = $(aobj).attr('href');
        $.getJSON(href, function (data) {
            $(aobj).hide();
            getPartVehicles();
            var notecount = 0;
            var vpID = 0;
            $(data).each(function (i, obj) {
                vpID = obj.ID;
                notecount += obj.Notes.length;
            });
            if (notecount == 0) {
                loadNotes(vpID);
            }
        });
    });

});

getPartVehicles = function () {
    var partid = $('#partID').val();
    $('#vehicleData').empty();
    $('#loadingVehicles').show();
    $.getJSON('/Product/GetPartVehicles', { partid: partid }, function (vData) {
        //console.log(vData);
        $('#loadingVehicles').hide();
        if (vData.length > 0) {
            $(vData).each(function (y, BaseVehicle) {
                var hasPart = (BaseVehicle.vehiclePart != null) ? true : false;
                var opt = '<li id="bv-' + BaseVehicle.ID + '">' + BaseVehicle.YearID + ' ' + BaseVehicle.Make.MakeName + ' ' + BaseVehicle.Model.ModelName + ((BaseVehicle.AAIABaseVehicleID != "") ? '<span class="vcdb">&#10004</span>' : '<span class="notvcdb">&times</span>');
                if (BaseVehicle.vehiclePart != null) {
                    opt += '<span class="tools">';
                    opt += '<a class="removeBV" href="/ACES/RemoveVehiclePart/' + BaseVehicle.vehiclePart.ID + '" title="Remove Vehicle Part Relationship">&times;</a>';
                    opt += ' | <a class="viewNotes" href="#" title="View Vehicle Part Notes" data-id="' + BaseVehicle.vehiclePart.ID + '">Notes</a>';

                    opt += '</span>';
                }
                opt += '<ul class="submodels">';
                $(BaseVehicle.Submodels).each(function (i, submodel) {
                    var hasPart = (submodel.vehiclePart != null) ? true : false;
                    opt += '<li id="bv' + BaseVehicle.ID + 's' + submodel.SubmodelID + '">' + submodel.submodel.SubmodelName.trim() + ((submodel.vcdb) ? '<span class="vcdb">&#10004</span>' : '<span class="notvcdb">&times</span>') + '<span class="tools">';
                    opt += '<span class="removalTools">';
                    if (submodel.vehiclePart != null) {
                        opt += '<a class="removeSubmodel" href="/ACES/RemoveVehiclePart/' + submodel.vehiclePart.ID + '" title="Remove Vehicle Part Relationship">&times;</a>';
                        opt += ' | <a class="viewNotes" href="#" title="View Vehicle Part Notes" data-id="' + submodel.vehiclePart.ID + '">Notes</a>';
                    }
                    opt += '</span>';
                    opt += '<a href="#" class="showConfig" title="Show / Hide Configurations">';
                    opt += '<span class="vehicleCount">' + submodel.vehicles.length + '</span><span class="arrow"></span>';
                    opt += '</a>';
                    opt += '</span><span class="clear"></span>';
                    opt += generatePartConfigTable(submodel);
                });
                opt += '</ul></li>';
                $('#vehicleData').append(opt);
                $('#vehicleData').find('div.configs').show();
            });
        } else {
            $('#vehicleData').append('<p>No Vehicles</p>');
        }
    });
};

generatePartConfigTable = function (submodel) {
    var showAdd = ($('#showAdd').val() != undefined && $('#showAdd').val() == 'true') ? true : false;
    var partid = ($('#partID').val() != undefined) ? $('#partID').val() : 0;
    var configTable = "";
    configTable += '<div class="configs"><table>';
    configTable += '<thead><tr>';
    configTable += '<th>VCDB</th>'
    $(submodel.configlist).each(function (z, config) {
        configTable += '<th>' + config.name + '</th>';
    });
    configTable += '<th></th>';
    configTable += '</tr></thead><tbody>';
    $(submodel.vehicles).each(function (x, vehicle) {
        var hasPart = (vehicle.vehiclePart != null) ? true : false;
        configTable += '<tr>';
        configTable += '<td>' + ((vehicle.vcdb) ? '<span class="vcdb">&#10004</span>' : '<span class="notvcdb">&times</span>') + '</td>';
        $(submodel.configlist).each(function (z, config) {
            configTable += '<td>';
            $(vehicle.configs).each(function (q, attr) {
                if (attr.ConfigAttributeType.name == config.name) {
                    configTable += attr.value;
                }
            });
            configTable += '</td>';
        });
        configTable += '<td>';
        if (vehicle.vehiclePart != null) {
            configTable += '<a class="removeConfig" href="/ACES/RemoveVehiclePart/' + vehicle.vehiclePart.ID + '" title="Remove Vehicle Part Relationship">&times;</a> | <a class="viewNotes" href="#" title="View Vehicle Part Notes" data-id="' + vehicle.vehiclePart.ID + '">Notes</a>';
        }
        configTable += '</td></tr>';
    });
    configTable += '</tbody></table></div>';
    return configTable;
};