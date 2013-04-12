var includedTable;
$(function () {
    includedTable = $('#includedParts').dataTable({ "bJQueryUI": true });

    $("#loading_area").fadeOut();
    $('#tableContainer').fadeIn();

    $('.remove').live('click', function () {
        var includedID = $(this).attr('id');
        var partID = $('#partID').val();
        var part = $(this).attr('title').substr(7, $(this).attr('title').length);
        if (partID > 0 && includedID > 0 && confirm('Are you sure you want to remove the relationship to ' + part + '?')) {
            $.getJSON('/Product/DeleteIncluded', { 'partID': partID, 'includedID': includedID }, function (response) {
                if (response.error == null) {
                    includedTable.fnDeleteRow($('#included\\:' + response.partID).get()[0]);
                    showMessage(response.shortDesc + ' removed.');
                } else {
                    showMessage(response.error);
                }
            });
        }
    });

    $(document).on('click', '#submitPart', function (e) {
        e.preventDefault();
        var bobj = $(this);
        var includedID = $('#addPart').val().trim();
        var quantity = $('#quantity').val();
        if (partID != "" && quantity != "") {
            var partID = $('#addPart').data('id');
            if(Number(quantity) <= 0) {
                showMessage('Quantity must be a number greater than 0');
                $('#quantity').focus();
                return;
            }
            $.getJSON('/Product/AddIncluded', { 'partID': partID, 'includedID': includedID, 'quantity': quantity }, function (response) {
                if (response.error == null) {
                    if ($('#included\\:' + response.includedID).get()[0] != null) {
                        includedTable.fnDeleteRow($('#included\\:' + response.includedID).get()[0]);
                    }
                    
                    // Add row to table
                    includedTable.fnAddData([
                            response.includedID,
                            response.quantity,
                            '<a href="javascript:void(0)" title="Remove ' + response.includedID + '" class="remove center removeincluded_' + response.includedID + '" id="' + response.includedID + '">Remove</a>'
                    ]);
                    $('.removeincluded_' + response.includedID).parent().parent().attr('id', 'included:' + response.includedID);
                    $('#addPart').attr('value', '');
                    $('#quantity').attr('value', '1');
                    showMessage(response.includedID + ' added.');
                } else {
                    showMessage(response.error);
                }
            });
        } else {
            $('#addPart').attr('value', '');
            showMessage("You must enter a part ID and have a quantity greater than 0.");
        }
    });

});