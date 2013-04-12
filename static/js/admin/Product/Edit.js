$(function () {

    $('#PartTypeID').on('blur', function () {
        if ($('#ACESPartTypeID').val() == "") {
            $('#PartTypeID').val('');
        } else {
            if ($('#PartTypeID').val().indexOf($('#ACESPartTypeID').val()) == -1) {
                $.getJSON('/ACES/GetPartTypeByID/' + $('#ACESPartTypeID').val(), function (resp) {
                    if (resp != null) {
                        $('#PartTypeID').attr('value', $.trim(resp.PartTerminologyName) + ' - ' + resp.PartTerminologyID);
                    } else {
                        $('#ACESPartTypeID').val()
                        $('#PartTypeID').val('');
                    }
                });
            }
        }
    });

    $('#PartTypeID').autocomplete({
        minLength: 2,
        source: function (request, response) {
            $.getJSON('/ACES/SearchPartTypes', { keyword: $('#PartTypeID').val() }, function (data) {
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
            $('#PartTypeID').val(ui.item.label);
            $('#ACESPartTypeID').val(ui.item.value);
        }
    });

});

