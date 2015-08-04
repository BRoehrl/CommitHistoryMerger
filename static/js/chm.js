function detailView() {
	document.getElementById("message").innerHTML = "Comment changed2.";
}

function getCommit(id) {
	$.getJSON('/commits/'+id, function(data){
		document.getElementById("message").innerHTML = data.Comment;
		document.getElementById("author").innerHTML = data.Author;
		document.getElementById("sha").innerHTML = data.Sha;
		document.getElementById("repository").innerHTML = data.Repo;
		document.getElementById("date").innerHTML = data.Time;
		document.getElementById("link").href = data.Link;
		document.getElementById("link").innerHTML = '../commit/'+data.Sha;
	});
	
}

