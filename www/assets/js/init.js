$(document).ready(function () {
    $('.sidenav').sidenav();
    $(".dropdown-trigger").dropdown({coverTrigger: false});
    $('.parallax').parallax();
    $('select').formSelect({dropdownOptions: {container: document.body}});
    $('.tooltipped').tooltip();
    $('.modal').modal();
    $('.datepicker').datepicker({format: 'dd-mm-yyyy'});
    $('.collapsible').collapsible();
});

function postAndRefresh(element, shouldConfirm) {
    if (shouldConfirm && !confirm("Are you sure you want to perform this action?")) return false;
    var $element = $(element);
    $.ajax({
        type: 'POST', url: $element.attr('href'),
        success: function () {
            window.location.reload(true);
        },
        error: function (e) {
            M.toast({html: e.responseText});
        }
    });
    return false;
}