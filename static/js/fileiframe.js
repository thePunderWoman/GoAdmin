function chooseFile() {
    $('#file-dialog').empty().html('<iframe id="fileViewerIFrame" src="/File/FrameIndex" frameborder="0" scrolling="yes" height="400" width="600"></iframe>');
    $("#file-dialog").dialog({
        autoOpen: false,
        height: 500,
        width: 700,
        title: "Choose a Photo",
        modal: true,
        buttons: {
            Cancel: function () {
                $(this).dialog("close");
            }
        },
        close: function () {
            $('#file-dialog').empty();
        }
    });

    $("#file-dialog").dialog("open");
}