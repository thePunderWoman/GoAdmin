var saveNewYear, saveNewMake, saveNewModel, saveNewSubmodel, loadYears, loadMakes, loadModels, loadSubmodel, updateMake, updateModel, updateSubmodel;

$(function () {
    $('#find').hide();
    $("#tabs").tabs();
    $('.addImg').show();

    $('#Addnonvcdb').on('submit', function (e) {
        e.preventDefault();
        var year = $('#nonyear').val();
        var makeID = $('#nonmake').val();
        var modelID = $('#nonmodel').val();
        var submodelID = $('#nonsubmodel').val();
        $.post('/ACES/AddNonVCDBVehicle', { nonyear: year, nonmake: makeID, nonmodel: modelID, nonsubmodel: submodelID }, function (data) {
            if (data.ID > 0) {
                showMessage("Vehicle Added successfully.");
            } else {
                showMessage("Error adding vehicle.");
            }
        },"json");
    });

    $(document).on('change', '#vcdbmake', function () {
        $('#find').hide();
        $('#vcdbmodel').html('<option value="">Select a Model</option>');
        $('#vcdbmodel').attr('disabled', 'disabled');
        var makeid = $(this).val();
        if (makeid != "") {
            $.getJSON('/ACES/GetVCDBModels/' + makeid, function (data) {
                $(data).each(function (i, model) {
                    $('#vcdbmodel').append('<option value="' + model.ModelID + '">' + $.trim(model.ModelName) + '</option>');
                });
                $('#vcdbmodel').removeAttr('disabled', 'disabled');
            })
        }
    });

    $(document).on('click', 'a.add', function (e) {
        e.preventDefault();
        var aobj = $(this);
        var href = $(aobj).attr('href');
        $.getJSON(href, function (data) {
            if (data.ID > 0) {
                $(aobj).hide();
                $(aobj).parent().append('<span class="added">Added</span>');
                $(aobj).parent().addClass('added');
            }
        })
    });
    

    $('#vcdbmodel').on('change', function (e) {
        if ($(this).val() == "") {
            $('#find').hide();
        } else {
            $('#find').show();
        }
    });

    $('#find').on('click', function () {
        var makeid = $('#vcdbmake').val();
        var modelid = $('#vcdbmodel').val();
        $('#vehicleData').empty();
        $('#loading').show();
        $.getJSON('/ACES/GetBaseVehicles', { makeid: makeid, modelid: modelid }, function (data) {
            $('#loading').hide();
            if (data.length > 0) {
                $(data).each(function (i, obj) {
                    var opt = '<li>' + obj.YearID + ' ' + obj.Make.MakeName + ' ' + obj.Model.ModelName + '<a href="/ACES/AddBaseVehicle/' + obj.BaseVehicleID + '" data-id="' + obj.BaseVehicleID + '" class="add">Add</a></li>';
                    $('#vehicleData').append(opt);
                });
            } else {
                $('#vehicleData').append('<p>No Unused Base Vehicles</p>');
            }
        });

    });

    $('#nonyear').change(function () {
        var yearID = $(this).val();
        if (yearID != "") {
            $('#delYear').fadeIn();
        } else {
            $('#delYear').fadeOut();
        }
    });

    $('#nonmake').change(function () {
        var makeID = $(this).val();
        var acesID = Number(makeID.split('|')[1]);
        if (acesID == 0) {
            $('#editMake').fadeIn();
            $('#delMake').fadeIn();
        } else {
            $('#editMake').fadeOut();
            $('#delMake').fadeOut();
        }
    });

    $('#nonmodel').change(function () {
        var modelID = $(this).val();
        var acesID = Number(modelID.split('|')[1]);
        if (acesID == 0) {
            $('#editModel').fadeIn();
            $('#delModel').fadeIn();
        } else {
            $('#editModel').fadeOut();
            $('#delModel').fadeOut();
        }
    });

    $('#nonsubmodel').change(function () {
        var modelID = $(this).val();
        var acesID = Number(modelID.split('|')[1]);
        if (acesID == 0) {
            $('#editSubmodel').fadeIn();
            $('#delSubmodel').fadeIn();
        } else {
            $('#editSubmodel').fadeOut();
            $('#delSubmodel').fadeOut();
        }
    });

    $('#addYear').click(function () {
        var html = '<input type="text" name="newYear" id="newYear" class="prompt_text" placeholder="Enter new year..." /><br />';
        $.prompt(html, {
            submit: saveNewYear,
            buttons: { Save: true }
        });
    });

    $('#delYear').live('click', function (e) {
        e.preventDefault();
        var yearID = $('#nonyear').val();

        if (yearID > 0 && confirm("Are you sure you want to remove this year?")) {
            $.getJSON('/ACES/RemoveYear', { 'year': yearID }, function (data) {
                if (data.success) {
                    loadYears();
                } else {
                    showMessage("There was a problem removing the year.")
                }
            });
        } else {
            if (yearID == 0) { // 
                showMessage('Invalid year.');
            }
        }
    });

    $('#addMake').click(function () {
        var html = '<input type="text" name="newMake" id="newMake" class="prompt_text" placeholder="Enter new make..." /><br />';
        $.prompt(html, {
            submit: saveNewMake,
            buttons: { Save: true }
        });
        $('#newMake').autocomplete({
            minLength: 1,
            source: function (request, response) {
                $.getJSON('/ACES/SearchMakes', { keyword: $('#newMake').val() }, function (data) {
                    response($.map(data, function (item) {
                        return {
                            label: item.label,
                            value: item.value,
                            id: item.ID
                        }
                    }));
                })
            },
            open: function () {
                $(this).removeClass("ui-corner-all").addClass("ui-corner-top");
            },
            close: function () {
                $(this).removeClass("ui-corner-top").addClass("ui-corner-all");
            },
            select: function (e, ui) {
                e.preventDefault();
                $('#newMake').val(ui.item.value);
            }
        });
    });

    $('#editMake').live('click', function () {
        var makeID = $('#nonmake').val();
        var mID = Number(makeID.split('|')[0]);
        if (mID > 0) {
            var html = '<input type="text" name="updateMake" id="updateMake" class="prompt_text" value="' + $('#nonmake option[value="' + makeID + '"]').text() + '" />';
            html += '<input type="hidden" name="makeID" id="makeID" value="' + mID + '" />';
            $.prompt(html, {
                submit: updateMake,
                buttons: { Save: true }
            });
        }
    });

    $('#delMake').live('click', function () {
        var makeID = $('#nonmake').val();
        var make = $('#nonmake option[value="' + makeID + '"]').text();
        if (makeID.length > 0 && confirm("Are you sure you want to remove " + make + "?")) {
            $.getJSON('/ACES/RemoveMake', { 'make': makeID }, function (data) {
                if (data.success) {
                    loadMakes();
                } else {
                    showMessage("There was a problem removing the make.")
                }
            });
        } else {
            if (makeID == "") { // 
                showMessage('Invalid make.');
            }
        }
    });

    $('#addModel').click(function () {
        var html = '<input type="text" name="newModel" id="newModel" class="prompt_text" placeholder="Enter new model..." />';
        $.prompt(html, {
            submit: saveNewModel,
            buttons: { Save: true }
        });
        $('#newModel').autocomplete({
            minLength: 1,
            source: function (request, response) {
                $.getJSON('/ACES/SearchModels', { keyword: $('#newModel').val() }, function (data) {
                    response($.map(data, function (item) {
                        return {
                            label: item.label,
                            value: item.value,
                            id: item.ID
                        }
                    }));
                })
            },
            open: function () {
                $(this).removeClass("ui-corner-all").addClass("ui-corner-top");
            },
            close: function () {
                $(this).removeClass("ui-corner-top").addClass("ui-corner-all");
            },
            select: function (e, ui) {
                e.preventDefault();
                $('#newModel').val(ui.item.value);
            }
        });
    });

    $('#editModel').live('click', function () {
        var modelID = $('#nonmodel').val();
        var mID = Number(modelID.split('|')[0]);
        if (mID > 0) {
            var html = '<input type="text" name="updateModel" id="updateModel" class="prompt_text" value="' + $('#nonmodel option[value="' + modelID + '"]').text() + '" />';
            html += '<input type="hidden" name="modelID" id="modelID" value="' + mID + '" />';
            $.prompt(html, {
                submit: updateModel,
                buttons: { Save: true }
            });
        }
    });

    $('#delModel').live('click', function () {
        var modelID = $('#nonmodel').val();
        var model = $('#nonmodel option[value="' + modelID + '"]').text();
        if (modelID.length > 0 && confirm("Are you sure you want to remove " + model + "?")) {
            $.getJSON('/ACES/RemoveModel', { 'model': modelID }, function (data) {
                if (data.success) {
                    loadModels();
                } else {
                    showMessage("There was a problem removing the model.")
                }
            });
        } else if (modelID == "") {
            showMessage('Invalid model.');
        }
    });

    $('#addSubmodel').click(function () {
        var html = '<input type="text" name="newSubmodel" id="newSubmodel" class="prompt_text" placeholder="Enter new submodel..." />';
        $.prompt(html, {
            submit: saveNewSubmodel,
            buttons: { Save: true }
        });
        $('#newSubmodel').autocomplete({
            minLength: 1,
            source: function (request, response) {
                $.getJSON('/ACES/SearchSubmodels', { keyword: $('#newSubmodel').val() }, function (data) {
                    response($.map(data, function (item) {
                        return {
                            label: item.label,
                            value: item.value,
                            id: item.ID
                        }
                    }));
                })
            },
            open: function () {
                $(this).removeClass("ui-corner-all").addClass("ui-corner-top");
            },
            close: function () {
                $(this).removeClass("ui-corner-top").addClass("ui-corner-all");
            },
            select: function (e, ui) {
                e.preventDefault();
                $('#newSubmodel').val(ui.item.value);
            }
        });
    });

    $('#editSubmodel').live('click', function () {
        var submodelID = $('#nonsubmodel').val();
        var sID = Number(submodelID.split('|')[0]);
        if (sID > 0) {
            var html = '<input type="text" name="updateSubmodel" id="updateSubmodel" class="prompt_text" value="' + $('#nonsubmodel option[value="' + submodelID + '"]').text() + '" />';
            html += '<input type="hidden" name="submodelID" id="submodelID" value="' + sID + '" />';
            $.prompt(html, {
                submit: updateSubmodel,
                buttons: { Save: true }
            });
        }
    });

    $('#delSubmodel').live('click', function () {
        var submodelID = $('#nonsubmodel').val();
        var submodel = $('#nonsubmodel option[value="' + submodelID + '"]').text();
        if (submodelID.length > 0 && confirm("Are you sure you want to remove " + submodel + "?")) {
            $.getJSON('/ACES/RemoveSubmodelByIDList', { 'submodel': submodelID }, function (data) {
                if (data.success) {
                    loadSubmodels();
                } else {
                    showMessage("There was a problem removing the submodel.")
                }
            });
        } else if (submodelID == "") {
            showMessage('Invalid submodel.');
        }
    });
});

saveNewYear = function (action, f, d, m) {
    var year = m.newYear;
    if (!isNaN(year) && year.length > 0 && year > 0) {
        $.getJSON('/ACES/AddYear', { 'year': year }, function (response) {
            loadYears();
        });
    } else {
        showMessage('Invalid year.');
    }
}

saveNewMake = function (action, f, d, m) {
    var make = m.newMake;
    if (make.trim().length > 0) {
        $.getJSON('/ACES/AddMake', { 'make': make}, function (response) {
            loadMakes();
        });
    } else {
        showMessage('Invalid make.');
    }
}

saveNewModel = function (action, f, d, m) {
    var model = m.newModel;
    if (model.trim().length > 0) {
        $.getJSON('/ACES/AddModel', { 'model': encodeURI(model) }, function (response) {
            loadModels();
        });
    } else {
        showMessage('Invalid model.');
    }
}

saveNewSubmodel = function (action, f, d, m) {
    var submodel = m.newSubmodel;
    if (submodel.trim().length > 0) {
        $.getJSON('/ACES/AddSubmodelByName', { 'submodel': encodeURI(submodel) }, function (response) {
            loadSubmodels();
        });
    } else {
        showMessage('Invalid model.');
    }
}

updateMake = function (action, f, d, m) {
    var make = m.updateMake;
    var makeID = m.makeID;
    if (make.length > 0) {
        $.getJSON('/ACES/UpdateMake', { 'id': makeID, 'name': make }, function (response) {
            loadMakes();
        });
    } else {
        showMessage('The make you entered is invalid.');
    }
}

updateModel = function (action, f, d, m) {
    var model = m.updateModel;
    var modelID = m.modelID;
    if (model.length > 0) {
        $.getJSON('/ACES/UpdateModel', { 'id': modelID, 'name': encodeURI(model) }, function (response) {
            loadModels();
        });
    } else {
        showMessage('The model you entered is invalid.');
    }
}

updateSubmodel = function (action, f, d, m) {
    var submodel = m.updateSubmodel;
    var submodelID = m.submodelID;
    if (submodel.length > 0) {
        $.getJSON('/ACES/UpdateSubmodel', { 'id': submodelID, 'name': encodeURI(submodel) }, function (response) {
            loadSubmodels();
        });
    } else {
        showMessage('The model you entered is invalid.');
    }
}

loadYears = function () {
    $.getJSON('/ACES/GetYears', function (years) {
        $('#nonyear').html('<option value="">- Select Year -</option>');
        $(years).each(function (i, year) {
            var opt = '<option value="' + year + '">' + year + '</option>'
            $('#nonyear').append(opt);
        });
        $('#nonyear').trigger('change');
    });
}

loadMakes = function () {
    $.getJSON('/ACES/GetAllMakes', function (makes) {
        $('#nonmake').html('<option value="">- Select Make -</option>');
        $.each(makes, function (i, make) {
            var new_option = '<option value="' + make.ID + '|' + ((make.AAIAID != null) ? make.AAIAID : 0) + '">' + make.name.trim() + '</option>';
            $('#nonmake').append(new_option);
        });
    });
}

loadModels = function () {
    $.getJSON('/ACES/GetAllModels', function (models) {
        $('#nonmodel').html('<option value="">- Select Model -</option>');
        $.each(models, function (i, model) {
            var new_option = '<option value="' + model.ID + '|' + ((model.AAIAID != null) ? model.AAIAID : 0) + '">' + model.name.trim() + '</option>';
            $('#nonmodel').append(new_option);
        });
    });
}

loadSubmodels = function () {
    $.getJSON('/ACES/GetAllSubmodels', function (submodels) {
        $('#nonsubmodel').html('<option value="">- Select Submodel -</option>');
        $.each(submodels, function (i, submodel) {
            var new_option = '<option value="' + submodel.ID + '|' + ((submodel.AAIAID != null) ? submodel.AAIAID : 0) + '">' + submodel.name.trim() + '</option>';
            $('#nonsubmodel').append(new_option);
        });
    });
}


/*String.prototype.trim = function() {
    return this.replace(/^\s+|\s+$/g,"");
}
String.prototype.ltrim = function() {
    return this.replace(/^\s+/,"");
}
String.prototype.rtrim = function() {
    return this.replace(/\s+$/,"");
}*/