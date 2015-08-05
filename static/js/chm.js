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
				dialog = document.getElementById('dialogAuthor')
				break;
			case 'Repository':
				dialog = document.getElementById('dialogAuthor')
				break;
			default:
				return;
		}
	dialog.showModal()

}

function closeDialog(filterBy) {
	switch(filterBy) {
			case 'Author':
				dialog = document.getElementById('dialogAuthor')
				break;
			case 'Repository':
				dialog = document.getElementById('dialogAuthor')
				break;
			default:
				return;
		}
	dialog.close()
}

