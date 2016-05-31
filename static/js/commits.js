function getCommit(id) {
  $.getJSON('./json/commits/' + id, function(data) {
    document.getElementById("message").innerHTML = data.Comment;
    document.getElementById("author").innerHTML = data.Author;
    document.getElementById("author").href = data.CreatorLink;
    document.getElementById("sha").innerHTML = data.Sha;
    document.getElementById("repository").innerHTML = data.Repo;
    document.getElementById("date").innerHTML = data.Time;
    document.getElementById("link").href = data.Link;
    document.getElementById("link").innerHTML = '../commit/' + data.Sha;
   });

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
