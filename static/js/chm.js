function showWaitDialog() {
  $('#pleaseWaitDialog').modal('show');
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

function refreshJWT() {
  $.get('./refresh_token');
}

refreshJWT();
setInterval(refreshJWT, 60*60*1000);
