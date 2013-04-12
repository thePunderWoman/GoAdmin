var tabIndex = 100;
$(function () {
    $('#partTable').dataTable({
        "bJQueryUI": true
    });

    $('.addImg:first').fadeIn();
    // Get the makes for this year and remove all other content for select boxes.
    $('#year').change(function () {
        var yearID = $(this).val();
        $('.delImg').fadeOut();
        $('.editImg').fadeOut();
        $('#editYear').fadeIn();
        if (yearID > 0) {
            $('#delYear').show();
            $.getJSON('/Vehicles/GetMakes', { 'yearID': yearID }, function (makes) {
                if (makes.length > 0) {
                    $('#make').html('<option value="0">- Select Make -</option>');
                    $('#model').html('<option value="0">- Select Model -</option>');
                    $('#style').html('<option value="0">- Select Style -</option>');
                }
                $.each(makes, function (i, make) {
                    var new_option = document.createElement('option');
                    $(new_option).val(make.makeID);
                    $(new_option).text(make.make1);
                    $('#make').append(new_option);
                });
            });
            $('.addImg').fadeOut();
            $('#col2').find('.addImg').fadeIn();

        }
    });

    // Get the models for this year and remove all other content for select boxes.
    $('#make').change(function () {
        var yearID = $('#year').val();
        var makeID = $(this).val();
        $('#delModel').fadeOut();
        $('#delStyle').fadeOut();
        $('#delMake').fadeIn();
        $('#editMake').fadeIn();
        $('#editYear').fadeIn();
        $('#editModel').fadeOut();
        $('#editStyle').fadeOut();
        $.getJSON('/Vehicles/GetModels', { 'yearID': yearID, 'makeID': makeID }, function (models) {
            if (models.length > 0) {
                $('#model').html('<option value="0">- Select Model -</option>');
                $('#style').html('<option value="0">- Select Style -</option>');
            }
            $.each(models, function (i, model) {
                var new_option = document.createElement('option');
                $(new_option).val(model.modelID);
                $(new_option).text(model.model1);
                $('#model').append(new_option);
            });
        });
        $('.addImg').fadeOut();
        $('#col3').find('.addImg').fadeIn();
    });

    // Get the styles for this year and remove all other content for select boxes.
    $('#model').change(function () {
        var yearID = $('#year').val();
        var makeID = $('#make').val();
        var modelID = $(this).val();
        $('#delStyle').fadeOut();
        $('#delModel').fadeIn();
        $('#editModel').fadeIn();
        $('#editStyle').fadeOut();
        $.getJSON('/Vehicles/GetStyles', { 'yearID': yearID, 'makeID': makeID, 'modelID': modelID }, function (styles) {
            if (styles.length > 0) {
                $('#style').html('<option value="0">- Select Style -</option>');
            }
            $.each(styles, function (i, style) {
                var new_option = document.createElement('option');
                $(new_option).val(style.styleID);
                $(new_option).text(style.style1);
                $('#style').append(new_option);
            });
        });
        $('.addImg').fadeOut();
        $('#col4').find('.addImg').fadeIn();
    });

    $('#style').change(function () {
        $('#editStyle').fadeIn();
        $('#delStyle').fadeIn();
    });

    $('#addYear').click(function () {
        var html = '<input type="text" name="newYear" id="newYear" class="prompt_text" placeholder="Enter new year..." /><br />';
        $.prompt(html, {
            submit: saveNewYear,
            buttons: { Save: true }
        });
    });

    $('#addMake').click(function () {
        var html = '<input type="text" name="newMake" id="newMake" class="prompt_text" placeholder="Enter new make..." /><br />';
        $.prompt(html, {
            submit: saveNewMake,
            buttons: { Save: true }
        });
    });

    $('#addModel').click(function () {
        var html = '<input type="text" name="newModel" id="newModel" class="prompt_text" placeholder="Enter new model..." />';
        $.prompt(html, {
            submit: saveNewModel,
            buttons: { Save: true }
        });
    });

    $('#addStyle').click(function () {
        var html = '<input type="text" name="newStyle" id="newStyle" class="prompt_text" placeholder="Enter new style..." />';
        $.prompt(html, {
            submit: saveNewStyle,
            buttons: { Save: true }
        });
    });

    $('#editYear').live('click', function () {
        var yearID = $('#year').val();
        if (yearID > 0) {
            var html = '<input type="text" name="updateYear" id="updateYear" class="prompt_text" value="' + $('#year option[value="' + yearID + '"]').text() + '" />';
            html += '<input type="hidden" name="yearID" id="yearID" value="' + yearID + '" />';
            $.prompt(html, {
                submit: updateYear,
                buttons: { Save: true }
            });
        }
    });

    $('#editMake').live('click', function () {
        var makeID = $('#make').val();
        if (makeID > 0) {
            var html = '<input type="text" name="updateMake" id="updateMake" class="prompt_text" value="' + $('#make option[value="' + makeID + '"]').text() + '" />';
            html += '<input type="hidden" name="makeID" id="makeID" value="' + makeID + '" />';
            $.prompt(html, {
                submit: updateMake,
                buttons: { Save: true }
            });
        }
    });


    $('#editModel').live('click', function () {
        var modelID = $('#model').val();
        if (modelID > 0) {
            var html = '<input type="text" name="updateModel" id="updateModel" class="prompt_text" value="' + $('#model option[value="' + modelID + '"]').text() + '" />';
            html += '<input type="hidden" name="modelID" id="modelID" value="' + modelID + '" />';
            $.prompt(html, {
                submit: updateModel,
                buttons: { Save: true }
            });
        }
    });

    $('#editStyle').live('click', function () {
        var styleID = $('#style').val();
        if (styleID > 0) {
            var html = '<input type="text" name="updateStyle" id="updateStyle" class="prompt_text" value="' + $('#style option[value="' + styleID + '"]').text() + '" />';
            html += '<input type="hidden" name="styleID" id="styleID" value="' + styleID + '" />';
            $.prompt(html, {
                submit: updateStyle,
                buttons: { Save: true }
            });
        }
    });

    $('#delYear').live('click', function () {
        var yearID = $('#year').val();

        if (yearID > 0 && confirm("Are you sure you want to remove this year? \r\n\r\nThis will remove any vehicles that are listed under this year!")) {
            $.getJSON('/Vehicles/DeleteYear', { 'yearID': yearID }, function (response) {
                response = $.trim(response);
                if (response.length > 0) {
                    showMessage(response);
                } else { // There was an error
                    $("#year option[value='" + yearID + "']").remove();
                    $('#year').val(0);
                    $('#delYear').fadeOut();
                    $('.addImg').fadeOut();
                    $('#addYear').fadeIn();
                    $('#make').html('<option value="0">- Select Make -</option>');
                    $('#model').html('<option value="0">- Select Model -</option>');
                    $('#style').html('<option value="0">- Select Style -</option>');
                    showMessage('Year removed.');
                }
            });
        } else {
            if (yearID == 0) { // 
                showMessage('Invalid year.');
            }
        }
    });

    $('#delMake').live('click', function () {
        var makeID = $('#make').val();
        var make = $('#make option[value="' + makeID + '"]').text();
        if (makeID > 0 && confirm("Are you sure you want to remove " + make + "? \r\n\r\nThis will remove any vehicles that are listed under this make!")) {
            $.getJSON('/Vehicles/DeleteMake', { 'makeID': makeID }, function (response) {
                response = $.trim(response);
                if (response.length > 0) {
                    showMessage(response);
                } else { // There was an error
                    $("#make option[value='" + makeID + "']").remove();
                    $('#make').val(0);
                    $('#delMake').fadeOut();
                    $('.addImg').fadeOut();
                    $('#addMake').fadeIn();
                    $('#model').html('<option value="0">- Select Model -</option>');
                    $('#style').html('<option value="0">- Select Style -</option>');
                    showMessage('Make removed.');
                }
            });
        } else {
            if (makeID == 0) { // 
                showMessage('Invalid make.');
            }
        }
    });

    $('#delModel').live('click', function () {
        var modelID = $('#model').val();
        var model = $('#model option[value="' + modelID + '"]').text();
        if (makeID > 0 && confirm("Are you sure you want to remove " + model + "? \r\n\r\nThis will remove any vehicles that are listed under this model!")) {
            $.getJSON('/Vehicles/DeleteModel', { 'modelID': modelID }, function (response) {
                response = $.trim(response);
                if (response.length > 0) {
                    showMessage(response);
                } else { // There was an error
                    $("#model option[value='" + modelID + "']").remove();
                    $('#model').val(0);
                    $('#delModel').fadeOut();
                    $('.addImg').fadeOut();
                    $('#addModel').fadeIn();
                    $('#style').html('<option value="0">- Select Style -</option>');
                    showMessage('Make removed.');
                }
            });
        } else if (modelID == 0) {
            showMessage('Invalid model.');
        }
    });

    $('#delStyle').live('click', function () {
        var styleID = $('#style').val();
        var style = $('#style option[value="' + styleID + '"]').text();
        if (styleID > 0 && confirm("Are you sure you want to remove " + style + "? \r\n\r\nThis will remove any vehicles that are listed under this style!")) {
            $.getJSON('/Vehicles/DeleteStyle', { 'styleID': styleID }, function (response) {
                response = $.trim(response);
                if (response.length > 0) {
                    showMessage(response);
                } else { // There was an error
                    $("#style option[value='" + styleID + "']").remove();
                    $('#style').val(0);
                    $('#delStyle').fadeOut();
                    $('.addImg').fadeOut();
                    $('#addStyle').fadeIn();
                    showMessage('Style removed.');
                }
            });
        } else if (styleID == 0) {
            showMessage('Invalid style.');
        }
    });
});

function saveNewYear(action,f,d,m) {
    var year = m.newYear;
    if (year.length > 0 && year > 0) {
        $.getJSON('/Vehicles/AddYear', { 'year': year }, function (response) {
            var yearID = response.yearID;
            if (yearID) {
                var new_year = document.createElement('option');
                $(new_year).val(yearID);
                $(new_year).text(response.year1);
                $('#year').append(new_year);
                $('#year').val(yearID);
                $('#delYear').fadeIn();
                $('#addMake').fadeIn();
                $('#year').trigger('change');
                showMessage(response.year1 + ' has been added.');
            } else {
                showMessage('Error: ' + response[0].error);
            }
        });
    } else {
        showMessage('Invalid year.');
    }
}

function saveNewMake(action, f, d, m) {
    var make = m.newMake;
    var yearID = $('#year').val();
    if (make.length > 0 && yearID > 0) {
        $.getJSON('/Vehicles/AddMake', { 'make': encodeURI(make), 'yearID': yearID }, function (response) {
            var makeID = response.makeID;
            if (makeID) {
                var new_make = document.createElement('option');
                $(new_make).val(makeID);
                $(new_make).text(response.make1);
                $('#make').append(new_make);
                $('#make').val(makeID);
                $('#delMake').fadeIn();
                $('#make').trigger('change');
                showMessage(response.make1 + ' has been added.');
            } else {
                showMessage('Error: ' + response[0].error);
            }
        });
    } else if (yearID == 0) {
        showMessage('Invalid year.');
    } else {
        showMessage('Error adding make.');
    }
}

function saveNewModel(action, f, d, m) {
    var model = m.newModel;
    var makeID = $('#make').val();
    var yearID = $('#year').val();
    if (model.length > 0 && makeID > 0 && yearID > 0) {
        $.getJSON('/Vehicles/AddModel', { 'model': encodeURI(model),'makeID':makeID }, function (response) {
            var modelID = response.modelID;
            if (modelID) {
                var new_model = document.createElement('option');
                $(new_model).val(modelID);
                $(new_model).text(response.model1);
                $('#model').append(new_model);
                $('#model').val(modelID);
                $('#delModel').fadeIn();
                $('#model').trigger('change');
                showMessage(response.model1 + ' has been added.');
            } else {
                showMessage('Error: ' + response[0].error);
            }
        });
    } else if (makeID == 0) {
        showMessage('Invalid make.');
    } else if (yearID == 0) {
        showMessage('Invalid year.');
    } else {
        showMessage('Error adding model.');
    }
}

function saveNewStyle(action, f, d, m) {
    var style = m.newStyle;
    var modelID = $('#model').val();
    if (style.length > 0 && modelID > 0) {
        $.getJSON('/Vehicles/AddStyle', { 'style': encodeURI(style), 'modelID': modelID }, function (response) {
            console.log(response);
            var styleID = response.styleID;
            if (styleID) {
                var new_style = document.createElement('option');
                $(new_style).val(styleID);
                $(new_style).text(response.style1);
                $('#style').append(new_style);
                $('#style').val(styleID);
                $('#delStyle').fadeIn();
                $('#style').trigger('change');
                showMessage(response.style1 + ' has been added.');
            } else {
                showMessage('Error: ' + response[0].error);
            }
        });
    } else if (style.length > 0) {
        showMessage('A model must be selected.');
    } else {
        showMessage('Invalid data.');
    }
}

function updateYear(action,f,d,m) {
    var year = m.updateYear;
    var yearID = m.yearID;
    if (year.length > 3) {
        $.getJSON('/Vehicles/EditYear', { 'year': year, 'yearID': yearID }, function (response) {
            var yearID = response.yearID;
            if (yearID) {
                $('#year option[value="' + yearID + '"]').text(response.year1);
                showMessage('Year updated.');
            } else {
                showMessage(response[0].error);
            }
        });
    } else {
        showMessage('The year you entered is invalid.');
    }
}

function updateMake(action, f, d, m) {
    var make = m.updateMake;
    var makeID = m.makeID;
    if (make.length > 0) {
        $.getJSON('/Vehicles/EditMake', { 'make': make, 'makeID': makeID }, function (response) {
            var makeID = response.makeID;
            if (makeID) {
                $('#make option[value="' + makeID + '"]').text(response.make1);
                showMessage('Make updated.');
            } else {
                showMessage(response[0].error);
            }
        });
    } else {
        showMessage('The make you entered is invalid.');
    }
}

function updateModel(action, f, d, m) {
    var model = m.updateModel;
    var modelID = m.modelID;
    if (model.length > 0) {
        $.getJSON('/Vehicles/EditModel', { 'model': model, 'modelID': modelID }, function (response) {
            var modelID = response.modelID;
            if (modelID) {
                $('#model option[value="' + modelID + '"]').text(response.model1);
                showMessage('Model updated.');
            } else {
                showMessage(response[0].error);
            }
        });
    } else {
        showMessage('The model you entered is invalid.');
    }
}

function updateStyle(action, f, d, m) {
    var style = m.updateStyle;
    var styleID = m.styleID;
    if (style.length > 0) {
        $.getJSON('/Vehicles/EditStyle', { 'style': style, 'styleID': styleID }, function (response) {
            var styleID = response.styleID;
            if (styleID) {
                $('#style option[value="' + styleID + '"]').text(response.style1);
                showMessage('Style updated.');
            } else {
                showMessage(response[0].error);
            }
        });
    } else {
        showMessage('The style you entered is invalid.');
    }
}