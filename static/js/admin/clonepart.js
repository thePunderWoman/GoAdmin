$(function () {
    $('#clonepart').click(clonePart);
});

function clonePart(event) {
    event.preventDefault();
    var partID, cloneformhtml;
    partID = $('#partID').val();
    cloneformhtml = '<label for="newPartID">New Part #<input type="text" id="newPartID" name="newPartID" value="" placeholder="Enter a new Part #..." /></label>';
    cloneformhtml += '<label for="newUPC">New UPC<input type="text" id="newUPC" name="newUPC" value="" placeholder="Enter a new UPC..." /></label>';
    cloneformhtml += '<p>What elements should be cloned?</p>';
    cloneformhtml += '<label for="cloneCats">Categories<input type="checkbox" id="cloneCats" value="true" /></label>';
    cloneformhtml += '<label for="cloneRelated">Related Parts<input type="checkbox" id="cloneRelated" value="true" /></label>';
    cloneformhtml += '<label for="cloneAttributes">Attributes<input type="checkbox" id="cloneAttributes" value="true" /></label>';
    cloneformhtml += '<label for="cloneContent">Content<input type="checkbox" id="cloneContent" value="true" /></label>';
    cloneformhtml += '<label for="cloneVehicles">Vehicles<input type="checkbox" id="cloneVehicles" value="true" /></label>';
    cloneformhtml += '<label for="clonePrices">Prices<input type="checkbox" id="clonePrices" value="true" /></label>';
    $("#clone-dialog").empty();
    $("#clone-dialog").append(cloneformhtml);
    $("#clone-dialog").dialog({
        modal: true,
        title: "Clone Part #" + partID + "?",
        buttons: {
            "Clone": function () {
                var newPartID, newUPC, categories, related, attributes, content, vehicles, prices;
                newPartID = $.trim($('#newPartID').val());
                newUPC = $.trim($('#newUPC').val());
                categories = $('#cloneCats').is(":checked");
                related = $('#cloneRelated').is(":checked");
                attributes = $('#cloneAttributes').is(":checked");
                content = $('#cloneContent').is(":checked");
                vehicles = $('#cloneVehicles').is(":checked");
                prices = $('#clonePrices').is(":checked");
                if (newPartID == "" || newUPC == "") {
                    $('#newPartID').addClass('error');
                    $('#newUPC').addClass('error');
                } else {
                    $(this).dialog("close");
                    $("#clone-dialog").empty();
                    $.getJSON('/Product/Clone', { 'partID': partID, 'newPartID': newPartID, 'upc': newUPC, 'categories': categories, 'relatedParts': related, 'attributes': attributes, 'content': content, 'vehicles': vehicles, 'prices': prices }, showMessages);
                }
            },
            Cancel: function () {
                $(this).dialog("close");
            }
        }
    });
};

function showMessages(messages) {
    $("#clone-dialog").append('<ul></ul>');
    $(messages).each(function (i, message) {
        $("#clone-dialog ul").append('<li>' + message + '</li>');
    });
    $("#clone-dialog").dialog({
        modal: true,
        buttons: {
            Ok: function () {
                $(this).dialog("close");
            }
        }
    });
}