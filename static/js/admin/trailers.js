$(function () {

    var trailerTable = $('table').dataTable({
        "bJQueryUI": true,
        "aaSorting": [[3, "asc"]]
    });

    $('#loading_area').slideUp();
    $('.dataTables_wrapper').fadeIn();
    $('table').fadeIn();


    $('.action').live('change', function () {
        var trailerID = $(this).attr('id');
        var catID = $('#categoryID').val();
        var action = $(this).val();

        switch (action) {
            case 'edit':
                // Redirect to the edit page for trailers
                window.location.href = '/Lifestyle/EditTrailer/' + trailerID;
                break;

            case 'delete':
                // Delete this category
                if (confirm("Are you sure you want to remove this trailer from this lifestyle?")) {
                    $.getJSON('/Lifestyle/RemoveTrailerFromLifestyle', { 'trailerid': trailerID, 'catID': catID }, function (response) {
                        if ($.trim(response).length == 0) {
                            trailerTable.fnDeleteRow($('#trailerRow\\:' + trailerID).get()[0]);
                            showMessage('Trailer has been removed.');
                        } else {
                            showMessage('' + response);
                        }
                    });
                }
                break;

            default:
                break;
        }
        $(this).val(0);
    });

    $('#AddTrailer').click(function (event) {
        event.preventDefault();
        var catID = $('#categoryID').val();
        $.getJSON('/Lifestyle/GetTrailersJSON', { 'catID': catID }, function (data) {
            $('#trailerselect').empty();
            if (data.length > 0) {
                $(data).each(function (i, obj) {
                    $('#trailerselect').append('<div class="trailer"><input name="trailers" type="checkbox" value="' + obj.trailerID + '" /><img src="' + ((obj.image != "") ? obj.image : "/Content/img/noimage.jpg") + '" alt="Image" /><p><strong>' + obj.name + '</strong><br />GTW: ' + obj.GTW + ' lbs<br />TW: ' + obj.TW + ' lbs</p><div class="clear"></div></div>');
                });
            } else {
                $('#trailerselect').append("<p>No Trailers To Add</p>");
            }

            $('#trailerselect').dialog({
                autoOpen: false,
                title: "Add a Trailer to Lifestyle",
                width: 480,
                height: 600,
                modal: true,
                buttons: {
                    "Done": function () {
                        var trailervals = "";
                        $.each($("input[name='trailers']:checked"), function () {
                            if(trailervals != "") trailervals += ","
                            trailervals += $(this).val();
                        });
                        $.getJSON('/Lifestyle/AddTrailersToLifestyle', { 'trailers': trailervals, 'catID': catID }, function (response) {
                            // Add to list
                            $(response).each(function (i, obj) {
                                trailerTable.fnAddData([
                                            obj.trailerID,
                                            obj.name,
                                            '<img class="trailerImage" src="' + ((obj.image == "") ? '/Content/img/noimage.jpg' : obj.image) + '" alt="Trailer Image" />',
                                            obj.TW,
                                            obj.GTW,
                                            '<select class="action" id="' + obj.trailerID + '"><option value="0">- Select Option -</option><option value="edit">Edit</option><option value="delete">Remove From Lifestyle</option></select>'
                                            ]);
                                $('#' + obj.trailerID).parent().parent().attr('id', 'trailerRow:' + obj.trailerID);
                            })
                            $('#trailerselect').dialog("close");
                        });
                    },
                    Cancel: function () {
                        $(this).dialog("close");
                    }
                },
                close: function () {
                    $('#trailerselect').empty();
                }
            });
            $('#trailerselect').dialog('open');
        });
    });

});