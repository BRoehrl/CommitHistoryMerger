function getCommit(id) {
  $.getJSON('./json/commits/' + id, function(data) {
    document.getElementById("message").innerHTML = data.Comment;
    document.getElementById("author").innerHTML = data.Author;
    document.getElementById("sha").innerHTML = data.Sha;
    document.getElementById("repository").innerHTML = data.Repo;
    document.getElementById("date").innerHTML = data.Time;
    document.getElementById("link").href = data.Link;
    document.getElementById("link").innerHTML = '../commit/' + data.Sha;
  });

}

function showWaitDialog() {
  $('#pleaseWaitDialog').modal('show');
}

function sendTag(tagType) {

  var tagBar = document.getElementById('tagBar');
  var tagInput;

  switch (tagType) {

    case 'Author':
      tagInput = $('#authorTag').parent().find('input.combobox').val();
      if (tagInput !== "") $('#tagBar').tagsinput('add', 'Author:' + tagInput);
      break;
    case 'Repository':
      tagInput = $('#repoTag').parent().find('input.combobox').val();
      if (tagInput !== "") $('#tagBar').tagsinput('add', 'Repo:' + tagInput);
      break;
    case 'Date':
      tagInput = document.getElementById('dateTag');
      if (tagInput !== "") $('#tagBar').tagsinput('add', 'Since:' + tagInput.value);
      break;
    default:
      return;
  }
  tagInput.value = "";

}

function saveProfile() {
  var nameInput = $('#profileName').parent().find('input.combobox').val();
  $.post('./config/save/' + nameInput);
  showWaitDialog();
  window.location.reload();
}

function loadProfile() {
  var selected = document.getElementById('selectedProfile');
  $.get('./config/load/' + selected.value);
  showWaitDialog();
  window.location.reload();
}

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
    // TODO Maybe show modal save dialog
  }

  $("#saveButton").toggle();
  $("#profileMenu").toggle();
}
