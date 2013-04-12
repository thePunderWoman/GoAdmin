$(function () {
    $(document).on('click', '#lines', function (e) {
        e.preventDefault();
        $('#sidebar').slideToggle('fast');
    });

    $(document).on('change', '#websiteID', function () {
        var websiteID = $(this).val();
        $.post("/Website/ChooseWebsite/" + websiteID, function (resp) {
            if (resp != "") {
                location.reload(true);
            }
        })
    });
});

String.prototype.trim = function () {
    return this.replace(/^\s+|\s+$/g, "");
}
String.prototype.ltrim = function () {
    return this.replace(/^\s+/, "");
}
String.prototype.rtrim = function () {
    return this.replace(/\s+$/, "");
}