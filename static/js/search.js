var mouseover_messagebox = false;
var resultbox_hoverstate = false;
var searchQueue = new Array();
$(function () {

    $('#docsSearch').submit(function (event) { event.preventDefault(); });

    $('#sidebar_border').mouseover(function (e) {
        if ($('#hide_arrow').length > 0) { $('#hide_arrow').remove(); }
        if ($('#show_arrow').length > 0) { $('#show_arrow').remove(); }
        if ($('#sidebar').css('display') != 'none') {
            $(this).css('background', '#686868').css('width', '4px');
            var xPos = parseInt(e.pageX) - 15;
            $('body').append('<img src="/Content/img/hide_arrow.gif" id="hide_arrow" style="z-index:20;position: fixed;top:' + e.pageY + 'px;left:' + xPos + 'px" />');
        } else {
            $(this).css('background', '#686868').css('width', '8px');
            var xPos = parseInt(e.pageX) + 15;
            $('body').append('<img src="/Content/img/show_arrow.gif" id="show_arrow" style="z-index:20;position: fixed;top:' + e.pageY + 'px;left:' + xPos + 'px" />');
        }
    });
    $('#sidebar_border').mouseout(function () {
        if ($('#sidebar').css('display') != 'none') {
            $(this).css('background', '#999').css('width', '3px');
        } else {
            $(this).css('background', '#999').css('width', '8px');
        }
        if ($('#hide_arrow').length > 0) { $('#hide_arrow').remove(); }
        if ($('#show_arrow').length > 0) { $('#show_arrow').remove(); }
    });

    $('#sidebar_border').click(function () {
        if ($('#sidebar').css('display') != 'none') {
            $('#sidebar').animate({ width: '0px' }).css('display', 'none');
            $(this).animate({ left: '0px', width: '8px' }, function () {
                $('#adminContent').animate({ width: '100%', margin: '75px 0px 0px 0px' });
            });
            $(this).css('width', '8px');
            $('#adminContent').css('margin-left', '0%');
        } else {
            $('#sidebar').animate({ width: '20%' }).css('display', 'block');
            $(this).animate({ left: '20%', width: '3px' }, function () {
                $('#adminContent').animate({ width: '80%', margin: '75px 0px 0px 20%' });
            });
            $('#adminContent').css('margin-left', '20%');
        }
    });

    $('#searchTerm').keyup(executeSearch);
    $('#searchTerm').focus(executeSearch);

    $('#searchResults').hover(function () { resultbox_hoverstate = true; }, hideResults);

    $('#message_box').mouseover(function () {
        mouseover_messagebox = true;
    }).mouseout(function () {
        mouseover_messagebox = false;
    });
});

function executeSearch() {
    var searchParam = $(this).val();
    if ($.trim(searchParam) != "") {
        var currentRequest = searchQueue.length + 1;
        searchQueue.push(currentRequest);
        $.getJSON('/Search/Search', { 'search_term': searchParam }, function (data) {
            if (resultbox_hoverstate == false && searchQueue[searchQueue.length - 1] == currentRequest) {
                var html = '';
                $.each(data, function (i, item) {
                    switch (item.type) {
                        case 'user':
                            html += '<li>';
                            html += '<a class="resultLink" href="/Users/Edit?user_id=' + item.id + '">' + item.term + '</a>';
                            html += '</li>';
                            break;
                        case 'category':
                            html += '<li>';
                            html += '<a class="resultLink" href="/Category/Edit?cat_id=' + item.id + '">' + item.term + '</a>';
                            html += '<li>';
                            break;
                        case 'part':
                            html += '<li>';
                            html += '<a class="resultLink" href="/Product/Edit?partID=' + item.id + '">Edit Part #' + item.term + '</a>';
                            html += '</li>';
                            break;
                        default:

                    }
                });
                $('#searchResults').find('ul').html(html);
                $('#searchResults').slideDown();
                searchQueue = [];
            }
        });
    } else {
        hideResults();
    }
}

function getMouseX(e) {
    return e.pageXOffset;
}

function hideResults() {
    resultbox_hoverstate = false;
    $('#searchResults').slideUp();
    $('#searchResults').html('<ul></ul>');
}

function showMessage(message) {
    $('#message_box').find('span').text(message);
    $('#message_box').fadeIn();
    setTimeout('hideMessage()', 5000);
}

function hideMessage() {
    if (!mouseover_messagebox) {
        $('#message_box').fadeOut();
    } else {
        setTimeout('hideMessage()', 800);
    }
}
