var testimonialTable = "";
$(function () {
    testimonialTable = $('#testimonialTable').dataTable({
        "bJQueryUI": true,
        "aoColumns": [
            null,
            null,
            null,
            null,
            null,
            { "sType": "date" },
            null,
            null
        ],
        "aaSorting": [[2, "asc"]]
    });

    $('.isApproved').live('click', function () {
        var testimonialID = $(this).attr('id').split(':')[1];
        if (testimonialID > 0) {
            $.getJSON('/Testimonial/Approve', { 'id': testimonialID }, function (response) {
                if (response.error == null) {
                    if (response == 1) {
                        showMessage("Testimonial Approved.");
                        testimonialTable.fnDeleteRow($('#' + testimonialID).parent().parent().get()[0]);
                    } else {
                        showMessage("Testimonial Unapproved.");
                        testimonialTable.fnDeleteRow($('#' + testimonialID).parent().parent().get()[0]);
                    }
                } else {
                    showMessage(response.error);
                }
            });
        }
    });

    $('.remove').live('click', function (event) {
        event.preventDefault();
        var testimonialID = $(this).attr('id');
        if (testimonialID > 0 && confirm('Are you sure you want to remove this testimonial?')) {
            var delete_testimonial = 0;
            $.get('/Testimonial/Remove', { 'id': testimonialID }, function (response) {
                if (response == "success") {
                    testimonialTable.fnDeleteRow($('#' + testimonialID).parent().parent().get()[0]);
                    showMessage('Testimonial removed.');
                } else {
                    showMessage(response);
                }
            });
        }

    });

});