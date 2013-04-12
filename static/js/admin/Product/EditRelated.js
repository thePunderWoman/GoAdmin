var relatedTable;
$(document).ready(function () {
    relatedTable = $('#relatedParts').dataTable({ "bJQueryUI": true });

    $(document).on('click', '.remove', function () {
        var relatedID = $(this).attr('id');
        var partID = $('#partID').val();
        var part = $(this).attr('title').substr(7, $(this).attr('title').length);
        console.log(partID)
        console.log(relatedID)
        if (partID > 0 && relatedID > 0 && confirm('Are you sure you want to remove the relationship to ' + part + '?')) {
            $.getJSON('/Product/DeleteRelated', { 'partID': partID, 'relatedID': relatedID }, function (response) {
                if (response.error == null) {
                    relatedTable.fnDeleteRow($('#related\\:' + response.partID).get()[0]);
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
        var partID = $('#addPart').data('id');
        var relatedID = $('#addPart').val().trim();
        if (partID > 0 && relatedID > 0 && confirm('Are you sure you want to make ' + relatedID + ' a related part?')) {
            // execute AJAX
            $.getJSON('/Product/AddRelated', { 'partID': partID, 'relatedID': relatedID }, function (response) {
                if (response.error == null) {
                    if ($('#related\\:' + response.partID).get()[0] != null) {
                        relatedTable.fnDeleteRow($('#related\\:' + response.partID).get()[0]);
                    }
                    // Add row to table
                    relatedTable.fnAddData([
                            response.partID,
                            response.shortDesc,
                            response.dateModified,
                            response.listPrice,
                            '<a href="javascript:void(0)" title="Remove ' + response.shortDesc + '" class="remove center" id="' + response.partID + '">Remove</a>'
                    ]);
                    $('#' + response.partID).parent().parent().attr('id', 'related:' + response.partID);
                    $('#addPart').attr('value', '');
                    showMessage(response.shortDesc + ' added.');
                } else {
                    showMessage(response.error);
                }
            });
        }
    });

});