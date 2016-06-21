var isLoading = false;
window.onload = function() {
    if ($('#tagBar').length) {
        var tags = sessionStorage.getItem('tags');
        if (tags !== null) {
            isLoading = true;
            $('#tagBar').val(tags);
            var tagArray = tags.split(",");
            for (var i = 0; i < tagArray.length; i++) {
                $('#tagBar').tagsinput('add', tagArray[i]);
            }
            isLoading = false;
            refreshQuery();
        } else {
            refreshQuery();
        }
    }
};

$('#tagBar').on('beforeItemAdd', function(event) {
    if (event.item.indexOf(':') == -1) {
      event.cancel = true;
      return;
    }
    // valid identifiers (author, repo, since, date)
    if('arsd'.indexOf(event.item[0].toLowerCase())  == -1) {
      event.cancel = true;
      return;
    }
});

$('#tagBar').on('itemAdded', function() {
    sessionStorage.setItem('tags', $('#tagBar').val());
    if (!isLoading) refreshQuery();
});

$('#tagBar').on('itemRemoved', function() {
    sessionStorage.setItem('tags', $('#tagBar').val());
    refreshQuery();
});

var compiledButton = _.template('<% _.each(button_data, function(bd) { %>\
  <button class="btn btn-lg btn-block btn-default" onclick="getCommit(\'<%= bd.ID %>\')" tstamp=<%= bd.NanoTime %> >\
    <div class="comCommentText">\
        <label>  <%= bd.Name %></label>\
    </div>\
    <div class="row comDateText">\
      <div class="col-xs-6 noColStyle pull-left">\
          <text>  <%= bd.DateString %></text>\
      </div>\
      <div class="col-xs-6 noColStyle pull-right noColRight hidden-xs">\
          <text>  <%= bd.Repository %></text>\
      </div>\
    </div>\
  </button> <% }); %>');

var loc = window.location;

function sortThat() {
    $("button.btn.btn-lg.btn-block.btn-default").sort(function(prev, next) {
        return $(next).attr("tstamp").localeCompare($(prev).attr("tstamp"));
    }).appendTo("#buttonList");
}


function deleteButtons() {
    document.getElementById('buttonList').innerHTML = "";
}


var searchBarValue = $('#searchBar').val().toLowerCase();
var authors = [];
var repos = [];
var dates = [];
var strDate = '';
var latestPage = 1;

function refreshQuery() {
    var i;
    authors = [];
    repos = [];
    dates = [];
    var items = $('#tagBar').tagsinput('items');
    var arrayLength = items.length;
    for (i = 0; i < arrayLength; i++) {
        var wholeTag = items[i];
        var type = wholeTag.substring(0, wholeTag.indexOf(":"));
        var tag = wholeTag.substring(wholeTag.indexOf(":") + 1);
        switch (type.toLowerCase()[0]) {
            case 'a': //author
                authors.push(tag);
                break;
            case 'r': //repo
                repos.push(tag);
                break;
            case 's': //since
            case 'd': //since (date)
                dates.push(tag);
                break;
            default:
                return;
        }
    }

    var query = './';
    if (authors.length > 0) {
        query = query + '&author=' + authors.join(';');
    }
    if (repos.length > 0) {
        query = query + '&repo=' + repos.join(';');
    }
    strDate = "";
    if (dates.length > 0) {
        var earliestDate = new Date(3000, 1, 1);
        var earliestString = '3000-01-01';
        for (i = 0; i < dates.length; i++) {
            var dateParts = dates[i].split("-");
            var date = new Date(dateParts[0], (dateParts[1] - 1), dateParts[2]);
            if (date <= earliestDate) {
                earliestDate = date;
                earliestString = dates[i];
            }
        }
        strDate = earliestString + 'T00:00:00Z';
        query = query + '&since=' + strDate;
    }
    deleteButtons();
    postQueryAndAddButtons({
        'page': 1
    }, function(request){
        var updatePage = request.getResponseHeader('Page-Has-Updates');
        if (updatePage === "true"){
          location.reload();
        }
    });

    query = query.replace('&', '?');
    var oldPath = "." + window.location.href.substring(window.location.href.lastIndexOf("/"));
    if (oldPath != encodeURI(query)) {
        window.history.pushState({
            "html": window.html,
            "pageTitle": window.pageTitle
        }, "", query);
        $('#pleaseWaitDialog').modal('show');
        window.location = query;
    }
}

function getCommit(id) {
    $.getJSON('./json/commits/' + id, function(data) {
        document.getElementById("message").innerHTML = data.Comment;
        var author = document.getElementById("author");
        author.innerHTML = data.Author;
        author.href = data.CreatorLink;
        if (data.CreatorLink === ""){
          author.className = "disabledLink";
        }else{
          author.className = "";
        }
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

function postQueryAndAddButtons(pagination, callbackFunction) {
    var params = {
        'commit': searchBarValue,
        'author': authors.join(';'),
        'repo': repos.join(';'),
        'date': strDate,
        'page': pagination.page,
        'perPage': 30,
    };
    $.ajax({
        type: 'POST',
        url: './commits',
        data: params,
        success: function(data, textStatus, request) {
            document.getElementById('buttonList').innerHTML += compiledButton({
                button_data: data
            });
            latestPage = pagination.page;
            if (typeof callbackFunction == 'function') callbackFunction(request);
        },
        dataType: "json",
        async: true
    });
}

function loadMore() {
    postQueryAndAddButtons({
        'page': latestPage + 1
    }, function() {
        $("#buttonList").bind('scroll', bindScroll);
    });
}

function bindScroll() {
    var vertical_margin = 30;
    var loadMore_treshhold = 500;
    if ($("#buttonList").scrollTop() + $("#buttonList").height() + vertical_margin > $("#buttonList")[0].scrollHeight - loadMore_treshhold) {
        $("#buttonList").unbind('scroll');
        loadMore();
    }
}

document.getElementById('searchBar').oninput = function() {
    searchBarValue = $('#searchBar').val().toLowerCase();
    deleteButtons();
    postQueryAndAddButtons({
        'page': 1
    });
};
$("#buttonList").scroll(bindScroll);
