var partID = 0;
var packageTable, savePackage, clearForm, showForm;

clearForm = (function () {
    var inch = $('#heightUnit option:contains(IN)').val();
    var pound = $('#heightUnit option:contains(LB)').val();
    var each = $('#heightUnit option:contains(EA)').val();
    $('#packageType').attr('value', '');
    $('#weight').attr('value', '');
    $('#height').attr('value', '');
    $('#width').attr('value', '');
    $('#length').attr('value', '');
    $('#quantity').attr('value', '1');
    $('#packageID').val(0);
    $('#heightUnit').val(inch)
    $('#lengthUnit').val(inch)
    $('#widthUnit').val(inch)
    $('#weightUnit').val(pound)
    $('#qtyUnit').val(each)
    $('.form_left').slideUp();
});

showForm = (function (packageID, weight, height, width, length, qty, packageUnit, dimensionUnit, weightUnit, type) {
    $('#packageID').attr('value', packageID);
    $('#weight').attr('value', weight);
    $('#height').attr('value', height);
    $('#width').attr('value', width);
    $('#length').attr('value', length);
    $('#quantity').attr('value', qty);
    $('#heightUnit').val(dimensionUnit)
    $('#lengthUnit').val(dimensionUnit)
    $('#widthUnit').val(dimensionUnit)
    $('#weightUnit').val(weightUnit)
    $('#qtyUnit').val(packageUnit)
    $('#packageType').val(type)
    $('.form_left').slideDown();
});

savePackage = (function (packageID, partID, weight, height, width, length, qty, weightUnit, heightUnit, widthUnit, lengthUnit, qtyUnit, type) {
    $.getJSON('/Product/SavePackage', { 'packageID': packageID, 'partID': partID, 'weight': weight, 'height': height, 'width': width, 'length': length, 'qty': qty, 'weightUnit': weightUnit, 'heightUnit': heightUnit, 'widthUnit': widthUnit, 'lengthUnit': lengthUnit, 'qtyUnit': qtyUnit, 'type': type }, function (response) {
        if (response.error == null) { // Success
            packageTable.fnAddData([
                    response.weight + " " + response.weightUnit.code,
                    response.height + " " + response.dimensionUnit.code,
                    response.width + " " + response.dimensionUnit.code,
                    response.length + " " + response.dimensionUnit.code,
                    response.quantity,
                    response.packageUnit.code,
                    response.PackageType.name,
                    '<a href="javascript:void(0)" class="edit" id="package_' + response.ID + '" data-id="' + response.ID + '">Edit</a> | <a href="javascript:void(0)" class="delete" data-id="' + response.ID + '">Delete</a>'
                    ]);
            clearForm();
        } else {
            showMessage(response.error);
        }
    });
});

$(function () {
    partID = $('#partID').val();
    packageTable = $('table').dataTable({ "bJQueryUI": true });

    $('#addPackage').on('click', function () {
        var inch = $('#heightUnit option:contains(IN)').val();
        var pound = $('#heightUnit option:contains(LB)').val();
        var each = $('#heightUnit option:contains(EA)').val();
        clearForm();
        showForm(0,"","","","",1,each,inch,pound,1);
    });

    $(document).on('click', '#btnReset', function () {
        var packageID = $('#packageID').val();
        if (packageID != 0) {
            $.getJSON('/Product/GetPackage', { 'packageID': packageID }, function (response) {
                if (response.error == null) { // Success
                    packageTable.fnAddData([
                    response.weight + " " + response.weightUnit.code,
                    response.height + " " + response.dimensionUnit.code,
                    response.width + " " + response.dimensionUnit.code,
                    response.length + " " + response.dimensionUnit.code,
                    response.quantity,
                    response.packageUnit.code,
                    response.PackageType.name,
                    '<a href="javascript:void(0)" class="edit" id="package_' + response.ID + '" data-id="' + response.ID + '">Edit</a> | <a href="javascript:void(0)" class="delete" data-id="' + response.packageID + '">Delete</a>'
                    ]);
                    clearForm();
                } else {
                    showMessage(response.error);
                }
            });
        }
        clearForm();
    });

    $(document).on('click', '#btnSave', function () {
        var packageID = $('#packageID').val();
        var type = $('#packageType').val();
        var weight = $('#weight').val();
        var height = $('#height').val();
        var width = $('#width').val();
        var length = $('#length').val();
        var qty = $('#quantity').val();
        var weightUnit = $('#weightUnit').val();
        var heightUnit = $('#heightUnit').val();
        var widthUnit = $('#widthUnit').val();
        var lengthUnit = $('#lengthUnit').val();
        var qtyUnit = $('#qtyUnit').val();
        savePackage(packageID, partID, weight, height, width, length, qty, weightUnit, heightUnit, widthUnit, lengthUnit, qtyUnit, type);
    });

    $(document).on('click', '.edit', function () {
        var row = $(this);
        var packageID = $(this).data('id');
        $.getJSON('/Product/GetPackage', { 'packageID': packageID }, function (response) {
            var weight = response.weight;
            var height = response.height;
            var width = response.width;
            var length = response.length;
            var qty = response.quantity;
            var packageUnit = response.packageUOM;
            var dimensionUnit = response.dimensionUOM;
            var weightUnit = response.weightUOM;
            var type = response.PackageType.ID;
            packageTable.fnDeleteRow($(row).parent().parent().get()[0]);
            showForm(packageID, weight, height, width, length, qty, packageUnit, dimensionUnit, weightUnit, type);
        });
    });

    $(document).on('click', '.delete', function () {
        var packageID = $(this).data('id');
        var table_row = $(this).parent().parent().get()[0];
        if (packageID > 0 && confirm('Are you sure you want to remove this package?')) {
            $.get('/Product/DeletePackage', { 'packageID': packageID }, function (response) {
                if (response.length == 0) {
                    packageTable.fnDeleteRow(table_row);
                } else {
                    showMessage(response);
                }
            });
            clearForm();
        }
    });
});