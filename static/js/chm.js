function getCommit(id) {
	$.getJSON('./json/commits/'+id, function(data){
		document.getElementById("message").innerHTML = data.Comment;
		document.getElementById("author").innerHTML = data.Author;
		document.getElementById("sha").innerHTML = data.Sha;
		document.getElementById("repository").innerHTML = data.Repo;
		document.getElementById("date").innerHTML = data.Time;
		document.getElementById("link").href = data.Link;
		document.getElementById("link").innerHTML = '../commit/'+data.Sha;
	});

}

function showWaitDialog(){
	$(pleaseWaitDialog).modal('show');
}

function sendTag(tagType){

	var tagBar = document.getElementById('tagBar');

	switch(tagType) {

			case 'Author':
				var tagInput = document.getElementById('authorTag')
				if (tagInput != "") $(tagBar).tagsinput('add', 'Author:' + tagInput.value)
				break;
			case 'Repository':
				var tagInput = document.getElementById('repoTag')
				if (tagInput != "") $(tagBar).tagsinput('add', 'Repo:' + tagInput.value)
				break;
			case 'Date':
				var tagInput = document.getElementById('dateTag')
				if (tagInput != "") $(tagBar).tagsinput('add', 'Since:' + tagInput.value)
				break;
			default:
				return;
		}
	tagInput.value = "";

}

function saveProfile(){
	var nameInput = document.getElementById('profileName')
	$.post('./config/save/'+ nameInput.value);
	closeDialog("Save")
	return false
}

function loadProfile(){
	var selected = document.getElementById('selectedProfile')
	$.get('./config/load/'+ selected.value);
	closeDialog("Load")
	window.location.href = './settings'
}

$("#settings :input").change(function() {
	$("#saveButton").toggle()
	$("#profileMenu").toggle()
	$("#settings").data("changed",true);
});

function saveSettings(){
	if ($("#settings").data("changed")) {
		var $form = $("#settings");
		var $inputs = $("#settings").find("input");
		var serializedData = $form.serialize();
		$.post('/settings', serializedData);
	}
	$("#saveButton").toggle()
	$("#profileMenu").toggle()
}
