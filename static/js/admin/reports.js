$(function() {
    $("#piesdatepicker").datepicker();
    $('#allpiesparts').on('click', function (e) {
        e.preventDefault();
        $('#piesdatepicker').attr('value', '');
        $('#piesreportform').submit();
    });
});