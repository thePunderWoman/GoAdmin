$(function () {
    $('.delete').live("click", function (event) {
        event.preventDefault();
        if (confirm("Are you sure you want to delete this?")) {
            window.location.href = event.currentTarget.href;
        }
    });
});