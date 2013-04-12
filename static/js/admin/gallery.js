$(function () {
    $('#sortgallery a:first').addClass('active').addClass('ascending');
    $('#sortgallery a').click(sortGallery);
    $('#loader').hide();
    $("#addgalleryform").submit(function (event) {
        event.preventDefault();
        var title = $('#galleryname').val();
        $.post('/File/AddGalleryAjax', { name: title, parentid: $('#parentid').val() }, function (response) {
            try {
                $('#galleryname').val('');
                var file_count = $('li.lifile').length;
                var gal_count = $('li.ligallery').length;
                var html = '<li id="gallery_' + response.fileGalleryID + '"><a class="gallery" href="/File/' + $('#location').val() + '/' + response.fileGalleryID + '"><span class="folder contextmenu"></span><span class="galleryname"><strong>' + response.name + '</strong> 0 subfolders; 0 files</span></a><ul class="menu"><li><a class="renamegallery" href="/File/RenameGallery/' + response.fileGalleryID + '" id="renamegallery_' + response.fileGalleryID + '">Rename</a></li><li><a class="deletegallery" href="/File/DeleteGallery/' + response.fileGalleryID + '" id="delgallery_' + response.fileGalleryID + '">Delete</a></li></ul></li>';
                if (file_count > 0) {
                    $('li.lifile:first').before(html);
                } else if (gal_count > 0) {
                    $('li.ligallery:last').after(html);
                } else {
                    $('#dropzone').after(html);
                }
            } catch (err) {
                console.log(response);
            };
        }, 'json');

    });
    $("html").click(function () { $("#galleries li ul.menu").hide(); });

    $(".deletegallery").live("click", function (event) {
        event.preventDefault();
        var idstr = $(this).attr('id').split("_")[1];
        if (confirm("Are you sure you want to delete this folder, all its files and subfolders? This cannot be undone.")) {
            $.post("/File/DeleteGalleryAjax?id=" + idstr, function (data) {
                if (data == "true") {
                    $('#gallery_' + idstr).remove();
                }
            }, "text");
        }
    });

    $(".deletefile").live("click", function (event) {
        event.preventDefault();
        var idstr = $(this).attr('id').split("_")[1];
        if (confirm("Are you sure you want to delete this file? This cannot be undone.")) {
            $.post("/File/DeleteFileAjax?id=" + idstr, function (data) {
                if (data == "true") {
                    $('#file_' + idstr).remove();
                }
            }, "text");
        }
    });

    $(".refreshfile").live("click", function (event) {
        event.preventDefault();
        var idstr = $(this).attr('id').split("_")[1];
        $.getJSON("/File/RefreshFileAjax?fileid=" + idstr, function (data) {
            var namespan = $('#file_' + idstr).find('span.filename');
            $(namespan).empty();
            var newinfo = '<strong>' + data.name + '</strong> path: <a href="' + data.path + '">link</a><br />';
            if (data.extension.FileType.fileType1.toLowerCase() == "image") {
                newinfo += 'dimensions: ' + data.height + ' x ' + data.width + '<br />';
            } else {
                newinfo += 'type: ' + data.extension.fileExt1 + '<br />';
            }
            var created = new Date(data.created);
            var ampm = "am";
            if (created.getHours() > 12) ampm = "pm";
            var datestr = (created.getMonth() + 1) + "/" + created.getDate() + "/" + created.getFullYear() + " " + ((created.getHours() > 12) ? (created.getHours() - 12) : created.getHours()) + ":" + created.getMinutes() + " " + ampm;
            newinfo += datestr + "<br />";
            if (data.size < 1024) {
                newinfo += data.size + ' Bytes';
            } else if (data.size >= 1024 && data.size < 1048576) {
                newinfo += Number(data.size / 1024).toFixed(2) + ' KB';
            } else if (data.size >= 1048576 && data.size < 1073741824) {
                newinfo += Number(data.size / 1048576).toFixed(2) + ' MB';
            } else {
                newinfo += Number(data.size / 1073741824).toFixed(2) + ' GB';
            }
            $(namespan).html(newinfo);
        });
    });

    $(".renamefile").live("click", function (event) {
        event.preventDefault();
        var idstr = $(this).attr('id').split("_")[1];
        $('#fileid').val(idstr);
        $('#galleryid').val(0);
        $('#newname').val($('#file_' + idstr).find('span.filename strong').text());
        $("#renameForm").dialog({
            autoOpen: false,
            height: 200,
            width: 350,
            title: "Rename File",
            modal: true,
            buttons: {
                "Rename": function () {
                    var name = $('#newname').val();
                    var galleryid = Number($('#galleryid').val());
                    var fileid = Number($('#fileid').val());
                    var bValid = true;

                    if ($.trim(name) == "") bValid = false;

                    if (bValid) {
                        $.post('/File/RenameAjax', { name: name, galleryid: galleryid, fileid: fileid }, function (data) {
                            $('#fileid').val(0);
                            $('#galleryid').val(0);
                            $('#newname').val('');
                            if (galleryid != 0) {
                                // replace folder name
                                $('#gallery_' + idstr).find('a.gallery span.galleryname strong').text(data.name);
                            } else {
                                // replace file name
                                $('#file_' + idstr).find('span.filename strong').text(data.name)
                            }
                        }, "json");
                        $(this).dialog("close");
                    }
                },
                Cancel: function () {
                    $(this).dialog("close");
                }
            },
            close: function () { }
        });
        $("#renameForm").dialog("open");
    });

    $(".renamegallery").live("click", function (event) {
        event.preventDefault();
        var idstr = $(this).attr('id').split("_")[1];
        $('#fileid').val(0);
        $('#galleryid').val(idstr);
        $('#newname').val($('#gallery_' + idstr).find('a.gallery span.galleryname strong').text());
        $("#renameForm").dialog({
            autoOpen: false,
            height: 200,
            width: 350,
            title: "Rename Folder",
            modal: true,
            buttons: {
                "Rename": function () {
                    var name = $('#newname').val();
                    var galleryid = Number($('#galleryid').val());
                    var fileid = Number($('#fileid').val());
                    var bValid = true;

                    if ($.trim(name) == "") bValid = false;

                    if (bValid) {
                        $.post('/File/RenameAjax', { name: name, galleryid: galleryid, fileid: fileid }, function (data) {
                            $('#fileid').val(0);
                            $('#galleryid').val(0);
                            $('#newname').val('');
                            if (galleryid != 0) {
                                // replace folder name
                                $('#gallery_' + idstr).find('a.gallery span.galleryname strong').text(data.name);
                            } else {
                                // replace file name
                                $('#file_' + idstr).find('span.filename strong').text(data.name)
                            }
                        }, "json");
                        $(this).dialog("close");
                    }
                },
                Cancel: function () {
                    $(this).dialog("close");
                }
            },
            close: function () { }
        });
        $("#renameForm").dialog("open");
    });



    $("#galleries li .contextmenu").live("contextmenu", function (event) {
        event.preventDefault();
        $("#galleries li ul.menu").hide();
        var posx = event.pageX - $(this).parent().parent().offset().left;
        var posy = event.pageY - $(this).parent().parent().offset().top;
        var menu = $(this).parent().parent().find('ul.menu');
        if (menu == null) $(this).parent().parent().find('ul.menu');
        menu.css('top', posy + 'px');
        menu.css('left', posx + 'px');
        menu.show();
    });

    if ($.browser.msie || !Modernizr.draganddrop) {
        $('#dropzone').remove();
        //$('#form_container').after('<p style="color:red">You\'re browser does not support drag and drop file uploading. Please upgrade to a different browser <br />(We recommend <a href="http://www.google.com/chrome" target="_blank">Google Chrome</a>.)</p>');
    } else {
        $('#dropzone').get(0).addEventListener('drop', upload, false);
        $('#dropzone').get(0).addEventListener('dragover', function (event) {
            event.preventDefault();
            $('#dropzone').css("background-color", "#ffc");
        }, false);
        $('#dropzone').get(0).addEventListener('dragleave', function (event) {
            event.preventDefault();
            $('#dropzone').css("background-color", "transparent");
        }, false);

    }

});

function upload(event) {
    event.stopPropagation();
    event.preventDefault();
    $('#dropzone').css("background-color", "transparent");
    var files = event.dataTransfer.files;
    for (var i = 0; i < files.length; i++) {
        uploadFile(files[i]);
    }
}

function uploadFile(file) {
    // Uploading - for Firefox, Google Chrome and Safari
    xhr = new XMLHttpRequest();


    // Add progress event listener
    xhr.upload.addEventListener("progress", function (evt) {
        if (evt.lengthComputable) {
            $('#loader').show();
            var loaded_pct = (evt.loaded / evt.total) * 100;
            $('#loader').attr('value', loaded_pct);
        } else {
            $('#loader').hide();
        }
    }, false);

    // Add event listener for the completed loading
    xhr.addEventListener("load", function (resp) {
        var response = resp.currentTarget.response;
        if (response != 'error') {
            var file = $.parseJSON(response);
            var type = file.extension.FileType.fileType1;
            var path = (type.toLowerCase() == "image") ? file.path : ((file.extension.fileExtIcon != "") ? file.extension.fileExtIcon : "/Content/img/file.png");
            var classstr = 'tall';
            if (file.height != 0) {
                classstr = ((file.width / file.height < 1.4) ? 'tall' : 'wide');
            }
            var file_count = $('li.lifile').length;
            var gal_count = $('li.ligallery').length;

            var html = '<li class="lifile" id="file_' + file.fileID + '"><span class="filebox"><img src="' + path + '" alt="' + file.name + '" class="' + classstr + ' contextmenu" /></span><span class="filename"><strong>' + file.name + '</strong> path: <a href="' + file.path + '" alt="direct link">link</a><br />' + ((type.toLowerCase() == 'image') ? 'dimensions: ' + file.height + ' x ' + file.width : 'type: ' + file.extension.fileExt1) + '<br />' + file.created + '<br />' + getFileSize(file.size) + '</span>';
            html += '<ul class="menu"><li><a class="renamefile" href="/File/RenameFile/' + file.fileID + '" id="renamefile_' + file.fileID + '">Rename</a></li><li><a class="deletefile" href="/File/DeleteFile/' + file.fileID + '" id="delfile_' + file.fileID + '">Delete</a></li></ul></li>';
            if (file_count > 0) {
                $('li.lifile:first').before(html);
            } else if (gal_count > 0) {
                $('li.ligallery:last').after(html);
            } else {
                $('#dropzone').after(html);
            }
        } else {
            showMessage('Error: Invalid file data.');
        }
        $('#loader').attr('value', '0');
        $('#loader').hide();
    }, false);

    var preserve = $('#preserve').is(':checked');
    xhr.open("post", "/File/AddFile", true);

    // Set appropriate headers
    xhr.setRequestHeader("Content-Type", "multipart/form-data");
    xhr.setRequestHeader("X-File-Name", file.name);
    xhr.setRequestHeader("X-File-Size", file.size);
    xhr.setRequestHeader("X-File-Type", file.type);
    xhr.setRequestHeader("X-Preserve-FileName", preserve);
    xhr.setRequestHeader("X-Gallery-ID", $('#parentid').val());

    // Send the file (doh)
    xhr.send(file);
}

function sortGallery(event) {
    event.preventDefault();
    var sortval = $(this).data("sort");
    var oldsort = $("#sortgallery a.active").data("sort");
    if (sortval == oldsort) {
        reverseSort();
    } else {
        // change and resort
        $("#sortgallery a.active").attr('class', '');
        $(this).addClass('active').addClass("ascending");
        var galleryid = $(this).data("galleryid");
        $.getJSON('/File/GetGalleryImagesJSON', { id: galleryid }, function (data) {
            var sortArray = new Array();
            var sortarray = new Array();
            $(data).each(function (i, obj) {
                sortarray.push(obj);
            });
            switch (sortval) {
                case "name":
                    sortarray.sort(sortByName);
                    break;
                case "date":
                    sortarray.sort(sortByDate);
                    break;
                case "type":
                    sortarray.sort(sortByType);
                    break;
                case "size":
                    sortarray.sort(sortBySize);
                    break;
            }
            var target = $('li.lifile:first').prev();
            $($(sortarray).get().reverse()).each(function (i,obj) {
                $(target).after($('#file_' + obj.fileID).detach());
            });
        });
    }
}

function reverseSort() {
    var direction = "ascending";
    if ($("#sortgallery a.active").hasClass('ascending')) direction = "descending";
    $("#sortgallery a.active").attr('class', 'active ' + direction);
    var files = $('li.lifile').get();
    var target = $('li.lifile:first').prev();
    $(files).each(function (i, obj) {
        $(target).after($(obj).detach());
    });
}

function sortBySize(a, b) {
    return a.size - b.size;
}

function sortByDate(a, b) {
    var datea = new Date(a.created);
    var dateb = new Date(b.created);
    return datea.getTime() - dateb.getTime();
}

function sortByType(a, b)
{
    // this sorts the array using the second element    
    return ((a.extension.fileExt1 < b.extension.fileExt1) ? -1 : ((a.extension.fileExt1 > b.extension.fileExt1) ? 1 : 0));
}

function sortByName(a, b) {
    // this sorts the array using the second element    
    return ((a.name < b.name) ? -1 : ((a.name > b.name) ? 1 : 0));
}

function getFileSize(size) {
    var bytestr = "";
    if(size < 1024) {
        bytestr = size + ' Bytes';
    } else if (size >= 1024 && size < 1048576) {
        bytestr = (size / 1024).toFixed(2) + ' KB';
    } else if (size >= 1048576 && size < 1073741824) {
        bytestr = (size / 1048576).toFixed(2) + ' MB';
    } else {
        bytestr = (size / 1073741824).toFixed(2) + ' GB';
    }
    return bytestr;
}