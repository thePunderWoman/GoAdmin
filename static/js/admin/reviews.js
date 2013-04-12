var reviewTable = "";
$(function () {
    reviewTable = $('#reviewTable').dataTable({
        "bJQueryUI": true,
        "aoColumns": [
            null,
            null,
            null,
            null,
            null,
            { "sType": "date" },
            null,
            null,
            null,
        ],
        "aaSorting": [[5, "asc"]]
    });

    $('.isApproved').live('click', function () {
        var reviewID = $(this).attr('id').split(':')[1];
        if (reviewID > 0) {
            $.getJSON('/Reviews/Approve', { 'id': reviewID }, function (response) {
                if (response.error == null) {
                    if (response == 1) {
                        showMessage("Review Approved.");
                        reviewTable.fnDeleteRow($('#' + reviewID).parent().parent().get()[0]);
                    } else {
                        showMessage("Review Unapproved.");
                    }
                } else {
                    showMessage(response.error);
                }
            });
        }
    });

    $('.action').live('change', function () {
        var action = $(this).val();
        var reviewID = $(this).attr('id');
        switch (action) {
            case 'view':
                if (reviewID > 0) {
                    $.getJSON('/Reviews/Get', { 'id': reviewID }, function (response) {
                        if (response.reviewID != null) {
                            $('#viewName').html('<strong>Name:</strong> ' + response.name);
                            $('#viewEmail').html('<strong>Email:</strong> ' + response.email);
                            $('#viewDate').html('<strong>Created On:</strong> ' + response.created);
                            $('#viewCustomer').html('<strong>Customer:</strong> ' + response.customer);
                            $('#viewRating').html('<strong>Rating:</strong> <div class="starrating"><span style="width: ' + ((response.rating / 5) * 100) + '%;"></div>');
                            $('#viewSubject').html('<strong>Subject:</strong> ' + response.subject);
                            $('#viewText').html('<strong>Review:</strong> ' + response.review_text);

                            $('#viewReview').dialog({
                                autoOpen: false,
                                title: "Review for Part #" + response.partID,
                                width: 800,
                                modal: true,
                                buttons: {
                                    "Approve": function () {
                                        $.getJSON('/Reviews/Approve', { 'id': response.reviewID }, function (approveresponse) {
                                            if (approveresponse.error == null) {
                                                if (approveresponse == 1) {
                                                    showMessage("Review Approved.");
                                                    reviewTable.fnDeleteRow($('#' + reviewID).parent().parent().get()[0]);
                                                } else {
                                                    showMessage("Review Unapproved.");
                                                }
                                            } else {
                                                showMessage(approveresponse.error);
                                            }
                                            $('#viewReview').dialog("close");
                                        });
                                    },
                                    Cancel: function () {
                                        $(this).dialog("close");
                                    }
                                },
                                close: function () {}
                            });
                            $('#viewReview').dialog('open');

                        } else if (response.error != null) {
                            showMessage(response.error);
                        } else {
                            showMessage('Something went wrong.');
                        }
                    });
                }
                break;

            case 'delete':
                if (reviewID > 0 && confirm('Are you sure you want to remove this review?')) {
                    var delete_review = 0;
                    var partID = $('#partID').val();
                    $.get('/Reviews/Remove', { 'id': reviewID }, function (response) {
                        if (response == "success") {
                            reviewTable.fnDeleteRow($('#' + reviewID).parent().parent().get()[0]);
                            showMessage('Review removed.');
                        } else {
                            showMessage(response);
                        }
                    });
                }
                break;

            default:
                break;
        }


        $(this).val(0);
    });

});