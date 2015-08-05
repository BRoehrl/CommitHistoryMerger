function detailView() {
	document.getElementById("message").innerHTML = "Comment changed2.";
}

function getCommit(id) {
	$.getJSON('/json/commits/'+id, function(data){
		document.getElementById("message").innerHTML = data.Comment;
		document.getElementById("author").innerHTML = data.Author;
		document.getElementById("sha").innerHTML = data.Sha;
		document.getElementById("repository").innerHTML = data.Repo;
		document.getElementById("date").innerHTML = data.Time;
		document.getElementById("link").href = data.Link;
		document.getElementById("link").innerHTML = '../commit/'+data.Sha;
	});
	
}

function openDialog(filterBy) {
	switch(filterBy) {
			case 'Author':
				var dialog = document.getElementById('dialogAuthor')
				break;
			case 'Repository':
				var dialog = document.getElementById('dialogRepo')
				break;
			case 'Date':
				var dialog = document.getElementById('dialogDate')
				break;
			default:
				return;
		}
	dialog.show()
}

function closeDialog(filterBy) {
	switch(filterBy) {
			case 'Author':
				var dialog = document.getElementById('dialogAuthor')
				break;
			case 'Repository':
				var dialog = document.getElementById('dialogRepo')
				break;
			case 'Date':
				var dialog = document.getElementById('dialogDate')
				break;
			default:
				return;
		}
	dialog.close()
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
	closeDialog(tagType)
}