$(function () {
    $('ul.imagesizes li:first a').addClass('active');
    $('ul.imagesize').hide();
    $('ul.imagesize:first').show();
    $('#size').val($('ul.imagesize:first').attr('id').split('_')[1])
    $('ul.imagesize:first').sortable({ handle: 'img', update: function (event, ui) { updateSort(event, ui) } });
    $('#import-images').click(function (event) {
        event.preventDefault();
        $(this).hide();
        var partid = $('#partID').val();
        $.post("/Product/ImportImages", { 'partid': partid }, function (response) {
            location.reload(true);
        }, "text");

    });
    $('a.deleteimage').live('click', function (event) {
        event.preventDefault();
        if (confirm("Are you sure you want to remove this image?")) {
            var idstr = $(this).attr('id').split('_')[1];
            $.post("/Product/DeleteImage", { imageid: idstr }, function (response) {
                if (response != "error") {
                    $("#partimage_" + idstr).fadeOut('fast', function () { $(this).remove(); });
                }
            }, "text");

        }
    });

    $('ul.imagesizes li a').click(function (event) {
        event.preventDefault();
        if (!$(this).hasClass('active')) {
            var o = $('ul.imagesizes li a.active').attr('href');
            $('ul.imagesizes li a.active').removeClass('active');
            $(this).addClass('active');
            var t = $(this).attr('href');
            $('ul.imagesize').hide();
            $('ul.imagesize' + t).fadeIn('fast');
            $('ul.imagesize' + t).sortable("destroy");
            $('ul.imagesize' + t).sortable({ handle: 'img', update: function (event, ui) { updateSort() } });
            $('#size').val(t.split('_')[1])
        }
    });

    $('#addimage').submit(function (event) {
        event.preventDefault();
        var file = $('#file').val();
        if ($.trim(file) != "") {
            var size = $('#size').val();
            var partid = $('#partID').val();
            $.post("/Product/AddImage", { size: size, partid: partid, file: file }, function (response) {
                if (response != "error") {
                    var data = $.parseJSON(response);
                    $('#size_' + size).append('<li class="partimage" id="partimage_' + data.imageID + '"><img src="' + data.path + '" alt="image-' + data.sort + '" /><span class="imagedetail">Dimensions: ' + data.width + ' x ' + data.height + '<br/><a href="#" class="deleteimage" id="deleteimage_' + data.imageID + '">Delete</a></span><span class="clear"></span></li>');
                    $('#file').val('');
                }
            }, "text");
        }
    })
});

function updateSort() {
    var t = $('ul.imagesizes li a.active').attr('href');
    var x = $('ul.imagesize' + t).sortable("serialize");
    $.post("/Product/updateSort?" + x);
}