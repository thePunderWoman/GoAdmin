var partID = 0;
var priceTable, savePrice, clearForm, showForm;

clearForm = (function () {
    $('#price_type').val(0);
    $('#price').val('');
    $('#priceID').val(0);
    $('#enforced').attr('checked', true);
    $('.form_left').slideUp();
});

showForm = (function (price_type, price, priceID, enforced) {
    if (price_type != null) { $('#price_type').val(price_type); }
    if (price != null) {
        price = price.replace('$', '');
        $('#price').val(price);
    }
    if (enforced == true) {
        $('#enforced').attr('checked', true);
    } else {
        $('#enforced').attr('checked', false);
    }
    if (priceID != null) { $('#priceID').val(priceID); }
    $('.form_left').slideDown();
});

savePrice = (function (price_type, price, priceID, enforced) {
    $.getJSON('/Product/SavePrice', { 'partID': partID, 'priceID': priceID, 'price': price, 'price_type': price_type, 'enforced': enforced }, function (response) {
        if (response.error == null) { // Success
            var exists = $('#price_' + response.priceID);
            if (exists.length > 0) {
                var table_row = $(exists).parent().parent().get()[0];
                priceTable.fnDeleteRow(table_row);
            }
            var enforcedString = "True";
            if (response.enforced.toString() == "false") {
                enforcedString = "False";
            }
         
            priceTable.fnAddData([
                    response.priceType,
                    '$' + response.price1.toString().replace('$', ''),
                    enforcedString,
                    response.dispDateModified,
                    '<a href="javascript:void(0)" class="edit" id="price_' + priceID + '" data-id="' + response.priceID + '">Edit</a> | <a href="javascript:void(0)" class="delete" data-id="' + response.priceID + '">Delete</a>'
                    ]);
            clearForm();
        } else {
            showMessage(response.error);
        }
    });
});

$(function () {
    partID = $('#partID').val();
    priceTable = $('table').dataTable({ "bJQueryUI": true });

    $('#addPrice').live('click', function () {
        clearForm();
        showForm();
    });

    $(document).on('click','#add_price_type', function () {
        var type = prompt("Enter the new price type.");
        if (type.length > 0) {
            $('#price_type').append("<option>" + type + "</option>");
            $('#price_type').val(type);
        }
    });

    $(document).on('click','#btnReset', function () {
        var priceID = $('#priceID').val();
        if (priceID != 0) {
            var price = $('#price').val();
            var type = $('#price_type').val();
            var enforced = "True";
            if ($('#enforced').is(":checked")) {
                enforced = "True";
            } else {
                enforced = "False";
            }
            priceTable.fnAddData([
                    type,
                    '$' + price.replace('$', ''),
                    enforced,
                    '<a href="javascript:void(0)" id="price_' + priceID + '" class="edit" data-id="' + priceID + '">Edit</a> | <a href="javascript:void(0)" class="delete" data-id="' + priceID + '">Delete</a>'
                ]);
        }
        clearForm();
    });

    $(document).on('click','#btnSave', function () {
        var priceID = $('#priceID').val();
        var price = $('#price').val();
        var price_type = $('#price_type').val();
        var enforced = true;
        if ($('#enforced').is(":checked")) {
            enforced = true;
        } else {
            enforced = false;
        }
        savePrice(price_type, price, priceID, enforced);
    });

    $(document).on('click','.edit', function () {
        var priceID = $(this).data('id');
        var enforced = true;
        
        if ($(this).parent().prev().prev().text().toString().toUpperCase() == "TRUE") {
            enforced = true;
        } else {
            enforced = false;
        }
        var price = $(this).parent().prev().prev().prev().text();
        var price_type = $(this).parent().prev().prev().prev().prev().text();
        
        priceTable.fnDeleteRow($(this).parent().parent().get()[0]);
        showForm(price_type, price, priceID, enforced);
    });

    $(document).on('click','.delete', function () {
        var priceID = $(this).data('id');
        var table_row = $(this).parent().parent().get()[0];
        if (priceID > 0 && confirm('Are you sure you want to remove this price record?')) {
            $.get('/Product/DeletePrice', { 'priceID': priceID }, function (response) {
                if (response.length == 0) {
                    priceTable.fnDeleteRow(table_row);
                } else {
                    showMessage(response);
                }
            });
            clearForm();
        }
    });
});