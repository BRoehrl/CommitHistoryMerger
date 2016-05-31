$("#settings :input").change(function() {
    if (!$("#settings").data("changed")) {
        $("#saveButton").toggle();
        $("#profileMenu").toggle();
        $("#settings").data("changed", true);
    }
});

function saveSettings() {
    if ($("#settings").data("changed")) {
        var $form = $("#settings");
        var $inputs = $("#settings").find("input");
        var serializedData = $form.serialize();
        $.post('/settings', serializedData);
        showSaveSuccessAlert();
        $("#settings").data("changed", false);
      }

    $("#saveButton").toggle();
    $("#profileMenu").toggle();
}

function showSaveSuccessAlert() {
    $("#save-success").alert();
    $("#save-success").fadeTo(2000, 500).slideUp(500);
}

$("#settings :input").keydown(function(e) {
    switch (e.which) {
        case 13: // enter
            this.blur();
            saveSettings();
            break;
        default:
            return; // exit this handler for other keys
    }
    e.preventDefault();
});
