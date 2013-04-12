/*
    Created By Jessica Janiuk
    Licensed under Creative Commons
    This file is designed to allow for multi-level drag and drop of items
    default structure looks like listed below:
    x = depth level.  starts at 0 for root items and adds with each child level
    y = current list item's unique id for reference in the DOM
    Depth starts at level 1

    <input type="hidden" id="children_0" value="comma separated list of the root level item ids" />
    <ul id="pages">
        <li class="level_x" id="item_y">
            <span class="handle">handle text / img here</span>
            <span class="title">title</span>
            <span id="meta_y">
                <input type="hidden" id="parent_y" value="y value for this item's parent" />
                <input type="hidden" id="children_y" value="comma separated list of y values for each of the children of this item in their order" />
                <input type="hidden" id="count_y" value="total number of this item's children" />
                <input type="hidden" id="sort_y" value="position in the sort of the parent of y" />
                <input type="hidden" id="depth_y" value="current level down, aka x value" />
            </span>
            <ul id="transport_y"></ul>
        </li>
    </ul>
*/

var currentdepth = 1;
$(document).ready(function () {
    bindStuff()
    $("a.remove").live('click', function (event) {
        event.preventDefault();
        removeItem($(this));
    });

    $("a.delete").live('click', function (event) {
        event.preventDefault();
        var idstr = $(this).attr('id').split('_')[1];
        $.getJSON("/Website/checkContent/" + idstr, function (data) {
            if (data.length == 0) {
                if (confirm("Are you sure you want to delete this content page?")) {
                    $.post("/Website/DeleteContent/" + idstr, function (data) {
                        if (data == "success") {
                            $('#page_' + idstr).fadeOut('fast', function () {
                                $('#page_' + idstr).remove();
                            })
                        }
                    })

                }
            } else {
                var alertmessage = "This content page is currently in " + data.length + " menus:";
                for (var i = 0; i < data.length; i++) {
                    alertmessage += "\n" + data[i].menuName
                }
                alertmessage += "\nYou must remove this page from those menus before it can be deleted."
                alert(alertmessage);
            }
        });
    });

    $('div.menutab').each(function () {
        if (!($(this).hasClass('active'))) {
            $(this).hide();
        }
    });

    $('#menutabs li').click(function () {
        if (!($(this).hasClass('active'))) {
            $('#menutabs li').removeClass('active');
            $(this).addClass('active');
            var target = $(this).attr('id').split('_')[1];
            $('div.menutab').hide();
            $('div#' + target).show();
        }
    });

});

var counter = 0;
function bindStuff() {
    $("ul#pages").sortable('destroy');
    $("ul#allpages").sortable('destroy');
    $("ul#pages").sortable({ handle: 'span.handle', placeholder: "page-placeholder", cursor: 'move', start: function (event, ui) { moveChildren(ui); setDepthValue(ui); }, sort: function (event, ui) { updateClasses(ui); }, beforeStop: function (event, ui) { setClasses(ui) }, stop: function (event, ui) { restoreChildren(ui); finishSort(ui); } }).disableSelection();
    $("ul#allpages").sortable({ handle: 'span.title', helper: 'original', cursor: 'move' }).disableSelection();
    $("li#tab_menulist").droppable({
        accept: "#allpages li",
        hoverClass: "dropper",
        tolerance: "pointer",
        drop: function (event, ui) {
            var obj = ui.draggable;
            contentid = obj.attr('id').split('_')[1];
            menuid = $('#menuid').val();
            $.post("/Website/AddContentToMenu?menuid=" + menuid + "&contentid=" + contentid, function (data) {
                var response = $.parseJSON(data);
                var published = "";
                if (response.published == true) published = " published";
                var liobj = "<li style=\"display:none;\" id=\"item_" + response.menuContentID + "\" class=\"level_1" + published + "\">" +
                    "<span class=\"handle\">↕</span> <span class=\"title\">" + response.pagetitle + "</span>" +
                    "<span class=\"controls\">" +
                        "<a href=\"/Website/SetPrimaryContent/" + response.contentID + "/" + menuid + "\"><img src=\"/Content/img/makeprimary.png\" alt=\"Make This Page the Primary Page\" title=\"Make This Page the Primary Page\" /></a> " +
                        "<a href=\"/Website/Content/Edit/" + response.contentID + "\"><img src=\"/Content/img/pencil.png\" alt=\"Edit Page\" title=\"Edit Page\" /></a> " +
                        " <a href=\"/Website/RemoveContent/" + response.menuContentID + "\" class=\"remove\" id=\"remove_" + response.menuContentID + "\"><img src=\"/Content/img/delete.png\" alt=\"Remove Page From Menu\" title=\"Remove Page From Menu\" /></a>" +
                    "</span>" +
                    "<span id=\"meta_" + response.menuContentID + "\">" +
                        "<input type=\"hidden\" id=\"parent_" + response.menuContentID + "\" value=\"" + ((response.parentID == null) ? 0 : response.parentID) + "\" />" +
                        "<input type=\"hidden\" id=\"children_" + response.menuContentID + "\" value=\"\" />" +
                        "<input type=\"hidden\" id=\"count_" + response.menuContentID + "\" value=\"0\" />" +
                        "<input type=\"hidden\" id=\"sort_" + response.menuContentID + "\" value=\"" + response.menuSort + "\" />" +
                        "<input type=\"hidden\" id=\"depth_" + response.menuContentID + "\" value=\"" + 1 + "\" />" +
                    "</span>" +
                    "<ul id=\"transport_" + response.menuContentID + "\"></ul>" +
                    "</li>";
                $('#pages').append(liobj);
                var childlist = $('#children_0').val();
                if (childlist != "") {
                    childlist += ",";
                }
                childlist += response.menuContentID;
                $('#children_0').val(childlist);
                $('#tab_menulist').trigger('click');
                $('#item_' + response.menuContentID).fadeIn();
            })
        }
    });
}

function hasChildren(idstr) {
    if($('#children_' + idstr).val() != '') {
        return true;
    }
    return false;
}

function isChild(idstr,parentid) {
    var childlist = $('#children_' + parentid).val().split(',');
    for(var i = 0;i < childlist.length; i++) {
        if(idstr == childlist[i]) return true;
    }
    return false;
}

function isLast(idstr) {
    var parentid = $('#parent_' + idstr).val();
    if(parentid != '') {
        var childlist = $('#children_' + parentid).val().split(',');
        var sortsize = childlist.length;
        if(Number($("#sort_" + idstr).val()) == sortsize) return true;
    } else {
        return true;
    }
    return false;
}

function finishSort(uiobj) {
    var idstr = $(uiobj.item).attr('id').split('_')[1];
    var prevparent = Number($('#parent_' + idstr).val());
    var serial = $("ul#pages").sortable("serialize").split('&');
    for (var i = 0; i < serial.length; i++) {
        var itemid = serial[i].split('=')[1];
        if (itemid == idstr) {
            var prev = getPreviousItem(uiobj.item,uiobj.item);
            var parent = getParentItem(uiobj.item,uiobj.item);
            if ($(parent).is('li')) {
                var parentid = Number($(parent).attr('id').split('_')[1]);
            } else {
                var parentid = 0;
            }

            // if changing parents, remove object id from original parent children string value
            if(prevparent == parentid) {
                //parents remain the same.  Adjust sorting accordingly.
                if (parentid == 0) {
                    // root element
                    var prevsame = getPreviousItemOnLevel(uiobj.item);
                    var previd = 0;
                    if (prevsame.length > 0) previd = Number($(prevsame).attr('id').split('_')[1]);
                    var childlist = $('#children_0').val().split(",")
                    var newsort = idstr;
                    if (previd == 0) {
                        // First in sort
                        var sortcount = 1;
                        $("#sort_" + idstr).val(sortcount);
                        for (var x = 0; x < childlist.length; x++) {
                            if (childlist[x] != idstr) {
                                sortcount++;
                                $("#sort_" + childlist[x]).val(sortcount);
                                newsort += "," + childlist[x];
                            }
                        }
                    } else {
                        var sortcount = 0;
                        for (var x = 0; x < childlist.length; x++) {
                            if (childlist[x] != idstr) {
                                sortcount++;
                                $("#sort_" + childlist[x]).val(sortcount);
                                if(x != 0) newsort += ","
                                newsort += childlist[x];
                                if (childlist[x] == previd) {
                                    sortcount++;
                                    $("#sort_" + idstr).val(sortcount);
                                    newsort += "," + idstr;
                                }

                            }
                        }
                    }
                    $('#children_' + parentid).val(newsort)
                } else {
                    if($(parent).attr('id') == $(prev).attr('id')) {
                        // first item in sort
                        var childlist = $('#children_' + parentid).val().split(",")
                        var newsort = idstr;
                        var sortcount = 1;
                        $("#sort_" + idstr).val(sortcount);

                        for(var x = 0; x < childlist.length;x++) {
                            if(childlist[x] != idstr) {
                                sortcount++;
                                $("#sort_" + childlist[x]).val(sortcount);
                                newsort += "," + childlist[x];
                            }
                        }
                        $('#children_' + parentid).val(newsort)
                    } else {
                        // middle of sort
                        var previd = $(prev).attr('id').split('_')[1];
                        var prevsort = Number($("#sort_" + previd).val());
                        var childlist = $('#children_' + parentid).val().split(",")
                        var newsort = "";
                        var sortcount = 0;
                        $("#sort_" + idstr).val(prevsort);
                        for(var x = 0; x < childlist.length;x++) {
                            if(childlist[x] != idstr) {
                                sortcount++;
                                $("#sort_" + childlist[x]).val(sortcount);
                                if(sortcount != 1) {
                                    newsort += ","    
                                } 
                                newsort += childlist[x];
                                if(childlist[x] == previd) {
                                    sortcount++;
                                    newsort += "," + idstr;
                                }
                            }
                        }
                        $('#children_' + parentid).val(newsort)
                    }
                }
            } else {
                // Remove id from previous parent
                var prevparentary = $('#children_' + prevparent).val().split(',');
                var prevparentstr = "";
                for(var x = 0; x < prevparentary.length; x++)
                {
                    if(Number(prevparentary[x]) != Number(idstr)) {
                        if(prevparentstr == '') {
                            prevparentstr = prevparentary[x];
                        } else {
                            prevparentstr = prevparentstr + "," + prevparentary[x];
                        }
                    }
                }
                $('#children_' + prevparent).val(prevparentstr);
                // set proper children count for previous parent
                setChildrenCounts(prevparent)

                if((parentid == 0 && !($(prev).is("li"))) || ($(parent).attr('id') == $(prev).attr('id'))) {
                    // item was dropped first in the sort list
                    // set sort = 1 and reset children values with this id as first
                    var childstring = $('#children_' + parentid).val();
                    $('#sort_' + idstr).val(1);
                    if(childstring == "" || childstring == idstr) {
                        // no children prior
                        $('#children_' + parentid).val(idstr);
                    } else {
                        // has children, loop through and reset sort vals
                        $('#children_' + parentid).val(idstr + "," + childstring);
                        var childlist = childstring.split(',');
                        for(var x = 0; x < childlist.length; x++) {
                            $('#sort_' + childlist[x]).val(Number($('#sort_' + childlist[x]).val()) + 1);
                        }
                    }
                } else {
                    /* item was dropped somewhere in the middle of the sort list
                    loop through and find prev id , update sort based on new position */
                    var childstring = $('#children_' + parentid).val();
                    var childlist = childstring.split(',');
                    var newchildren = "";
                    var sort = 0;
                    var prev = getPreviousItemOnLevel(uiobj.item);
                    var previd = Number($(prev).attr('id').split('_')[1]);
                    for(var x = 0; x < childlist.length; x++) {
                        sort++;
                        if(x == 0) {
                            newchildren = childlist[x];
                        } else {
                            newchildren = newchildren + "," + childlist[x];
                        }
                        if(Number(childlist[x]) == previd) {
                            newchildren = newchildren + "," + idstr;
                            sort++;
                            $('#sort_' + idstr).val(sort);
                        } else {
                            $('#sort_' + childlist[x]).val(sort);
                        }
                    }
                    $('#children_' + parentid).val(newchildren);
                }
            }
            $('#parent_' + idstr).val(parentid);
            if(parentid != 0) setChildrenCounts(parentid);
        }
    }
    saveItems()
}

/* Methods happen on start of sort to ensure child items move properly with parent item */
function moveChildren(uiobj) {
    var idstr = $(uiobj.helper).attr('id').split('_')[1];
    if (hasChildren(idstr)) {
        var transport = $('#transport_' + idstr);
        getChildren(idstr,transport);
    }
    $(uiobj.placeholder).css('height',$(transport).height() + $(uiobj.helper).height() + 'px');
}

function getChildren(idstr,transport) {
    var childlist = $('#children_' + idstr).val().split(',');
    // loop and get all items in between the start and end
    for(var i = 0; i < childlist.length; i++) {
        var liobj = $('#item_' + childlist[i])
        $(transport).append(liobj);
        if(hasChildren(childlist[i])) {
            getChildren(childlist[i],transport)
        }
    }
}

/* Happens on drop of parent item to move child items back into place */
function restoreChildren(uiobj) {
    var idstr = $(uiobj.item).attr('id').split('_')[1];
    var elements = $('#transport_' + idstr).children();
    $(elements).insertAfter($(uiobj.item))
}

/* Happens during drag to check where the item is currently hovering and ensure that it can't drop too deep or too shallow */
function getMaxLevel(prevobj) {
    var maxlevel = 0;
    if($(prevobj).is('li')) {
        var previd = $(prevobj).attr('id').split('_')[1];
        var nextlevel = Number($('#depth_' + previd).val());
        var maxlevel = nextlevel + 1;
    }
    return maxlevel;
}

function getMinLevel(prevobj,currentobj) {
    var level = 1;
    var currid = Number($(currentobj).attr('id').split('_')[1]);
    var currparent = Number($('#parent_' + currid).val());
    var islast = false;
    if(Number($('#sort_' + currid).val()) == Number($('#count_' + currid).val())) {
        var islast = true;
    }

    if($(prevobj).is('li'))
    {
        var previd = Number($(prevobj).attr('id').split('_')[1]);
        var prevsort = Number($('#sort_' + previd).val());
        var prevchildren = $('#children_' + previd).val().split(',');
        var prevparent = Number($('#parent_' + previd).val());
        var prevdepth = Number($('#depth_' + previd).val());
        var prevcount = Number($('#count_' + previd).val());
        if((prevsort < prevcount) && !(islast && (currparent == prevparent))){
            // previous item is inside a set of children.  Min level should equal this item's level.
            if(currentdepth == prevdepth && prevchildren.length < 2) {
                // ok to match level
                level = Number($('#depth_' + previd).val());
            } else {
                // previous is the parent of a set of children. Check for how many children.
                if(prevchildren.length == 1) {
                    if(prevchildren[0] != currid && prevchildren[0] != "") {
                        level = Number($('#depth_' + previd).val()) + 1;
                    } else {
                        level = Number($('#depth_' + previd).val());
                    }
                } else {
                    level = Number($('#depth_' + previd).val()) + 1;
                }
            }
        } else {
            // previous item is the end of the set of children.
             if(prevchildren[0] != "" && prevchildren[0] != currid) {
                // moved to an item that had only one child
                level = Number($('#depth_' + previd).val()) + 1;
            } else {

                var parentobj = getParentItem(prevobj,prevobj);
                if($(parentobj).is('li')) {
                    var parentid = Number($(parentobj).attr('id').split('_')[1]);
                } else {
                    var parentid = 0;
                }
                while(parentid != 0) {
                    var parentid = Number($(parentobj).attr('id').split('_')[1]);
                    var parentparent = Number($('#parent_' + parentid).val());
                    var parentsort = Number($('#sort_' + parentid).val());
                    var parentcount = Number($('#count_' + parentid).val());
                    var parentdepth = Number($('#depth_' + parentid).val());
                    var parentchildren = $('#children_' + parentid).val().split(',');
                    if((parentcount > parentsort) && !(islast && (parentparent == currparent))) {
                        level = parentdepth;
                        parentid = 0;
                    }
                    if(parentid != 0) {
                        var parentobj = getParentItem(parentobj,parentobj);
                        if(!($(parentobj).is('li'))){
                            parentid = 0;
                        }
                    }
                }
            }
        }
    }
    return level;
}

/* Gets the parent item of elements after drop */
function getParentItem(obj,helper,depth) {
    var objid = Number($(helper).attr('id').split('_')[1]);
    var prevobj = getPreviousItem(obj, helper);
    var parentobj = ''
    if($(prevobj).is("li") && $(prevobj).attr('id') != undefined) {
        var previd = Number($(prevobj).attr('id').split('_')[1]);
        var prevlevel = Number($('#depth_' + previd).val());
        var objlevel = Number($('#depth_' + objid).val());
        if(depth != undefined) objlevel = Number(depth);
        if(prevlevel < objlevel) {
            // previous object is current obj's parent
            parentobj = prevobj;
        } else if(prevlevel == objlevel) {
            // prev object is on same level.  Get Prev obj parent to get correct parent.
            var parentobj = $('#item_' + $('#parent_' + previd).val());
        } else {
            // prev object is the child of a previous object.  Get parents until one level above.
            for(var i = prevlevel; i >= objlevel; i--)
            {
                var previd = Number($(prevobj).attr('id').split('_')[1]);
                var prevobj = $('#item_' + $('#parent_' + previd).val());
            }
            parentobj = prevobj;
        }
    }
    return parentobj;
}

function getPreviousItemOnLevel(obj) {
    var idstr = $(obj).attr('id').split('_')[1];
    var targetlevel = Number($("#depth_" + idstr).val());
    var prevobj = getPreviousItem(obj, obj);
    if (prevobj.length > 0) {
        var previdstr = $(prevobj).attr('id').split('_')[1];
        var prevlevel = Number($("#depth_" + previdstr).val());
        if (prevlevel == targetlevel) {
            return prevobj;
        } else {
            var nextobj = "";
            while (targetlevel < prevlevel) {
                var parentstr = 'item_' + $("#parent_" + previdstr).val();
                var nextobj = $('#' + parentstr);
                previdstr = Number($(nextobj).attr('id').split('_')[1]);
                prevlevel = Number($("#depth_" + previdstr).val());
            }
            return nextobj;
        }
    }
    return "";
}

/* Gets item above what was just grabbed and / or dropped */
function getPreviousItem(obj,helper) {
    var prevobj = $(obj).prev();
    if($(prevobj).attr('id') == $(helper).attr('id')) {
        var prevobj = $(obj).prev().prev()
    }
    return prevobj;
}

/* Happens during drag to set placeholder class values */
function updateClasses(obj) {
    var idstr = $(obj.helper).attr('id').split('_')[1];
    var currentlevel = Number($("#depth_" + idstr).val());
    var ph = obj.placeholder;
    var previousobj = getPreviousItem(ph,obj.helper)
    var parentobj = getParentItem(ph,obj.helper,currentdepth)
    var minlevel = 1;
    var maxlevel = getMaxLevel(previousobj);
    var currentpix = currentlevel * 20;
    var pixdif = Number(obj.helper.css('left').split('p')[0]) + currentpix;
    var targetlevel = Math.floor(pixdif/20) ;
    if(targetlevel >= maxlevel) targetlevel = maxlevel;
    var minlevel = getMinLevel(previousobj,obj.helper);
    if(targetlevel < minlevel) targetlevel = minlevel;
    $(ph).attr('class', 'page-placeholder level_' + targetlevel);
    currentdepth = targetlevel;
}

function setDepthValue(obj) {
    var idstr = $(obj.helper).attr('id').split('_')[1];
    currentdepth = Number($("#depth_" + idstr).val());
}

/* Happens right as you drop to set all final class values */
function setClasses(obj) {
    var idstr = $(obj.helper).attr('id').split('_')[1];
    var classstring = $(obj.placeholder).attr('class').split('page-placeholder ')[1];
    var currentlevel = Number($("#depth_" + idstr).val());
    var newlevel = Number(classstring.split('_')[1]);
    if(currentlevel != newlevel) {
        $("#depth_" + idstr).val(newlevel);
        if($(obj.helper).hasClass('published')) {
            $(obj.helper).attr('class','level_' + newlevel + ' published');
        } else {
            $(obj.helper).attr('class','level_' + newlevel);
        }
        if(hasChildren(idstr)) {
            setChildDepth(idstr)
        }
    }
}

/* Loops through all child elements to set proper depth and class */
function setChildDepth(idstr) {
    var childlist = $('#children_' + idstr).val().split(',');
    var parentdepth = Number($("#depth_" + idstr).val());
    // loop and get all items in between the start and end
    for(var i = 0; i < childlist.length; i++) {
        var currentdepth = $("#depth_" + childlist[i]).val();
        $("#depth_" + childlist[i]).val(parentdepth + 1);
        $("#item_" + childlist[i]).removeClass('level_' + currentdepth);
        $("#item_" + childlist[i]).addClass('level_' + (parentdepth + 1));
        if(hasChildren(childlist[i])) {
            setChildDepth(childlist[i]);
        }
    }
}

function setChildrenCounts(idstr) {
    var childlist = $('#children_' + idstr).val().split(',');
    for(var i = 0; i < childlist.length; i++) {
        $('#count_' + childlist[i]).val(childlist.length);
    }
}

function saveItems() {
    var datastring = "";
    var count = 0
    $('ul#pages li').each(function () {
        count++;
        var idstr = $(this).attr('id').split('_')[1];
        if (count == 1) {
            datastring += "?";
        } else {
            datastring += "&";
        }
        datastring += "page[]=" + idstr + "-" + $("#parent_" + idstr).val() + "-" + $("#sort_" + idstr).val()
    });

    $.post("/Website/MenuSort/" + $('#menuid').val() + datastring)
}

function removeItem(obj) {
    var idstr = $(obj).attr("id").split('_')[1];
    var item = $('#item_' + idstr);
    var level = $('#depth_' + idstr).val();
    var parent = $('#parent_' + idstr).val();
    if (hasChildren(idstr)) {
        setChildData($('#children_' + idstr).val(), parent);
        // replace parent childlist with new children
        var parentchildren = $('#children_' + parent).val().split(',');
        var newchildren = ""
        for (var i = 0; i < parentchildren.length; i++) {
            if(i != 0) newchildren += ",";
            if (Number(parentchildren[i]) == Number(idstr)) {
                newchildren += $('#children_' + idstr).val();
            } else {
                newchildren += parentchildren[i];
            }
        }
        $('#children_' + parent).val(newchildren)
    }
    // update sort for all items on this level
    $.post("/Website/RemoveContentAjax/" + idstr)
    $(item).fadeOut('fast', function () {
        $(item).remove();
        sort();
    });
}

function sort() {
    var level = 1
    var parent = 0;
    var list = "";
    var sort = 0;
    $('#pages li').each(function () {
        var id = Number($(this).attr('id').split('_')[1]);
        if (parent == Number($('#parent_' + id).val())) {
            // root element
            sort++;
            if ($(this).hasClass('published')) {
                $(this).attr('class', 'level_' + level + ' published');
            } else {
                $(this).attr('class', 'level_' + level);
            }
            $('#depth_' + id).val(level);
            $('#sort_' + id).val(sort);
            if (hasChildren(id)) {
                var children = $('#children_' + id).val();
                sortChildren(children, (level + 1), id);
            }
            if (list == "") {
                list = id;
            } else {
                list += "," + id;
            }
        }
    })
    $('#children_' + parent).val(list);
}

function sortChildren(childlist, level, parent) {
    var children = childlist.split(',');
    var sort = 0;
    var list = "";
    for(var i = 0; i < children.length; i++) {
        var item = $('#item_' + children[i]);
        // check if item exists
        if(item.length > 0) {
            sort++;
            if ($(item).hasClass('published')) {
                $(item).attr('class', 'level_' + level + ' published');
            } else {
                $(item).attr('class', 'level_' + level);
            }
            $('#depth_' + children[i]).val(level);
            $('#sort_' + +children[i]).val(sort);
            if(hasChildren(children[i])) {
                var childvals = $('#children_' + children[i]).val();
                sortChildren(childvals,(level + 1),children[i]);
            }
            if(list == "") {
                list = children[i];
            } else {
                list += "," + children[i];
            }
        }
    }
    $('#children_' + parent).val(list);
}

function setChildData(childlist,parent)  {
    var children = childlist.split(',');
    var sort = 0;
    for(var i = 0; i < children.length; i++) {
        $('#parent_' + children[i]).val(parent);
    }
}