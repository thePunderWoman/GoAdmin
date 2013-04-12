$(document).ready(function () {
    var partCatTable = $('table').dataTable({ "bJQueryUI": true });

    // Add this part to a category
    $('#addCat').click(function () {
        var catID = $('#catID').val();
        var cat = $('#catID option[value="' + catID + '"]').text();
        var partID = $('#partID').val();
        $.get('/Product/AddCategory', { 'catID': catID, 'partID': partID }, function (response) {
            if (response == "") {
                partCatTable.fnAddData([
                        cat,
                        '<a href="javascript:void(0)" class="remove" id="' + catID + '">Remove</a>'
                    ]);
                showMessage("Category added.");
                $('#catID').val(0);
            } else {
                showMessage(response);
            }
        });
    });

    // Remove this part from a category
    $('.remove').live('click', function () {
        var catID = $(this).attr('id');
        var cat = $(this).parent().prev().text();
        var part = $('#shortDesc').val();
        var partID = $('#partID').val();
        var tableRow = $(this).parent().parent().get()[0];
        if (partID > 0 && catID > 0 && confirm("Are you sure you want to remove " + part + " from " + cat + "?")) {
            $.get('/Product/DeleteCategory', { 'catID': catID, 'partID': partID }, function (response) {
                if ($.trim(response) == "") {
                    partCatTable.fnDeleteRow(tableRow);
                    showMessage("Category removed.");
                } else {
                    showMessage(response);
                }
            });
        }
    });

});