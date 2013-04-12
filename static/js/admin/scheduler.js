var deleteTask, runTask, changeFreq;

deleteTask = function (name, id, href) {
    if (confirm('Are you sure you want to delete task "' + name + '"?')) {
        $.post('/Scheduler/DeleteTask', { 'id': id }, function () {
            $('#task_' + id).fadeOut('fast',function () { $('#task_' + id).remove(); });
        });
    }
};

runTask = function (id) {
    $('#task_' + id + ' img.running').show();
    $.post('/Scheduler/RunTask', { 'id': id }, function () {
        $('#task_' + id + ' img.running').hide();
    });
};

changeFreq = function () {
    var freq = $('#runfrequency').val();
    switch (freq) {
        case "interval":
            $('#rundaylabel').hide();
            $('#runtimelabel').hide();
            $('#intervallabel').show();
            break;
        case "daily":
            $('#rundaylabel').hide();
            $('#runtimelabel').show();
            $('#intervallabel').hide();
            break;
        case "weekly":
            $('#rundaylabel').show();
            $('#runtimelabel').show();
            $('#intervallabel').hide();
            break;
        case "monthly":
            $('#rundaylabel').show();
            $('#runtimelabel').show();
            $('#intervallabel').hide();
            break;
    }
}

$(function () {
    changeFreq();
    $('#runtime').timepicker({ ampm: true });
    $('#runday').datepicker();
    $('#runfrequency').on('change', changeFreq);
    $(document).on('click', '.delete', function (e) {
        event.preventDefault();
        var id, name, href;
        id = $(this).data('id');
        name = $(this).data('name');
        href = $(this).attr('href');
        deleteTask(name, id, href);
    });

    $(document).on('click', '.runtask', function (e) {
        event.preventDefault();
        var id;
        id = $(this).data('id');
        runTask(id);
    });
});