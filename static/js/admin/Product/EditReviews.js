var reviewTable = "";
$(document).ready(function () {
    reviewTable = $('#reviewTable').dataTable({
        "bJQueryUI": true,
        "aoColumns": [
            null,
            null,
            null,
            null,
            null,
            {"sType": "date"},
            null,
            null,
            null,
        ],
        "aaSorting": [[5, "desc"]]
    });

    $('#addReview img').click(function () {
        $('#name').val('');
        $('#email').val('');
        $('#review_text').val('');
        $('#subject').val('');
        $('#rating').val(0);
        $('#reviewID').val(0);
        $('#reviewForm').slideDown();
        $(this).fadeOut();
    });

    $('#btnReset').click(function () {
        resetAddForm();
        return false;
    });

    $('#btnSubmit').click(function () {
        var name = $('#name').val();
        var email = $('#email').val();
        var subject = $('#subject').val();
        var review_text = $('#review_text').val();
        var rating = $('#rating').val();
        var partID = $('#partID').val();
        var reviewID = $('#reviewID').val();
        if (reviewID == 0) {
            addReview(partID, rating, subject, review_text, name, email);
        } else {
            updateReview(partID, rating, subject, review_text, name, email, reviewID);
        }
    });

    $('.isApproved').live('click', function () {
        var reviewID = $(this).attr('id').split(':')[1];
        if (reviewID > 0) {
            $.getJSON('/Reviews/Approve', { 'id': reviewID }, function (response) {
                if (response.error == null) {
                    if (response == 1) {
                        showMessage("Review Approved.");
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
                                close: function () { }
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

            case 'edit':
                if (reviewID > 0) {
                    $.getJSON('/Reviews/Get', { 'id': reviewID }, function (response) {
                        if (response.reviewID != null) {
                            $('#name').val(response.name);
                            $('#email').val(response.email);
                            $('#rating').val(response.rating);
                            $('#subject').val(response.subject);
                            $('#review_text').val(response.review_text);
                            $('#reviewID').val(response.reviewID);
                            $('#reviewForm').slideDown();
                            $('#addReview img').fadeOut();
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
                        if (response == "") {
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

function resetAddForm() {
    $('#name').val('');
    $('#email').val('');
    $('#review_text').val('');
    $('#subject').val('');
    $('#rating').val(0);
    $('#reviewID').val(0);
    $('#reviewForm').slideUp();
    $('#addReview img').fadeIn();
}

function addReview(partID, rating, subject, review_text, name, email) {
    if (partID > 0 && review_text.length > 0 && rating > 0) {
        $.getJSON('/Product/AddReview', { 'partID': partID, 'rating': rating, 'subject': subject, 'review_text': review_text, 'name': name, 'email': email }, function (response) {
            if (response.error == null) {
                reviewTable.fnAddData([
                            response.reviewID,
                            response.partID,
                            '<div class="starrating"><span style="width: ' + ((Number(response.rating) / 5) * 100) + '%;"></span></div>',
                            response.customer,
                            response.name,
                            response.created,
                            response.subject,
                            '<input type="checkbox" name="isApproved" id="isApproved:' + response.reviewID + '" value="1" />',
                            '<select name="action" class="action" id="' + response.reviewID + '"><option value="0">- Select Action -</option><option value="view">View</option><option value="edit">Edit</option><option value="delete">Delete</option></select>'
                        ]);
                resetAddForm();
                showMessage('Review added.');
            } else {
                showMessage(response.error);
            }
        });
    } else {
        if (review.length == 0) {
            showMessage('Review is required.');
        } else if (rating == 0) {
            showMessage('You must enter a rating.');
        } else {
            showMessage('Invalid Part #');
        }
    }
}

function updateReview(partID, rating, subject, review_text, name, email, reviewID) {
    if (partID > 0 && review_text.length > 0 && rating > 0 && reviewID > 0) {
        reviewTable.fnDeleteRow($('#' + reviewID).parent().parent().get()[0]);
        $.getJSON('/Product/AddReview', { 'partID': partID, 'rating': rating, 'subject': subject, 'review_text': review_text, 'name': name, 'email': email, 'reviewID': reviewID }, function (response) {
            if (response.error == null) {
                reviewTable.fnAddData([
                            response.reviewID,
                            response.partID,
                            '<div class="starrating"><span style="width: ' + ((Number(response.rating) / 5) * 100) + '%;"></span></div>',
                            response.customer,
                            response.name,
                            response.created,
                            response.subject,
                            '<input type="checkbox" name="isApproved" id="isApproved:' + response.reviewID + '" value="1" />',
                            '<select name="action" class="action" id="' + response.reviewID + '"><option value="0">- Select Action -</option><option value="view">View</option><option value="edit">Edit</option><option value="delete">Delete</option></select>'
                        ]);
                resetAddForm();
                showMessage('Review saved.');
            } else {
                showMessage(response.error);
            }
        });
    } else {
        if (review_text.length == 0) {
            showMessage('Review Text is required.');
        } else if (rating == 0) {
            showMessage('You must enter a rating.');
        } else {
            showMessage('Invalid Part #');
        }
    }
}