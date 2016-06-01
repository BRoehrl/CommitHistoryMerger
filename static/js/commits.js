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

function loadMore() {
    console.log("More loaded");
    //$("#buttonList").append("<button class=\"btn btn-lg btn-block btn-default\" onclick=\"getCommit('3d19fac3d87f2e51328fa972446e3d21e12f2043')\"><div class=\"comCommentText\"><label>  Updating develop poms back to pre merge state</label>				</div>				<div class=\"row comDateText\"> <div class=\"col-xs-6 noColStyle pull-left\"> <text>  30 May 16 </text></div> <div class=\"col-xs-6 noColStyle pull-right noColRight hidden-xs\"><text>ingrid-portal/develop</text></div></div></button>");
    $("#buttonList").bind('scroll', bindScroll);
}

function bindScroll() {
    var vertical_margin = 30;
    var loadMore_treshhold = 2000;
    if ($("#buttonList").scrollTop() + $("#buttonList").height() + vertical_margin > $("#buttonList")[0].scrollHeight - loadMore_treshhold) {
        $("#buttonList").unbind('scroll');
        loadMore();
    }
}

function searchCommits(){
  var searchText = $( '#searchBar' ).val().toLowerCase();
  $("#buttonList").children('button').each(function ( index ){
    var labelText = $( this ).find('div > label').text().toLowerCase();
    if (labelText.indexOf(searchText) === -1){
      $( this ).hide();
    }else {
      $( this ).show();
    }
  });
}

document.getElementById('searchBar').oninput = searchCommits;
$("#buttonList").scroll(bindScroll);
