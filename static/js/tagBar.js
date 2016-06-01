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
        }else{
          refreshQuery();
        }
    }
};

$('#tagBar').on('beforeItemAdd', function(event) {
    if (event.item.indexOf(':') == -1) {
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
var serversocket = new WebSocket("ws://" + loc.host + "/socket");
serversocket.onopen = function() {
    //refreshQuery();
};

function sortThat() {
    $("button.btn.btn-lg.btn-block.btn-default").sort(function(prev, next) {
        return $(next).attr("tstamp").localeCompare($(prev).attr("tstamp"));
    }).appendTo("#buttonList");
}


serversocket.onmessage = function(e) {
    console.log(e);
    if (loc.pathname !== "/") return;
    document.getElementById('buttonList').innerHTML += compiledButton({
        button_data: $.parseJSON(e.data)
    });
};

function deleteButtons() {
    document.getElementById('buttonList').innerHTML = "";
}

var authors = [];
var repos = [];
var dates = [];
var strDate = '';
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
        switch (type.toLowerCase()) {
            case 'author':
                authors.push(tag);
                break;
            case 'repo':
                repos.push(tag);
                break;
            case 'since':
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

    postQueryAndSetButtons();

    query = query.replace('&', '?');
    var oldPath = "." + window.location.href.substring(window.location.href.lastIndexOf("/"));
    if (oldPath != encodeURI(query)) {
        window.history.pushState({
            "html": window.html,
            "pageTitle": window.pageTitle
        }, "", query);
        //$('#pleaseWaitDialog').modal('show');
        //window.location = query;
    }
}

function postQueryAndSetButtons(){
  var params = {
      'author': authors.join(';'),
      'repo': repos.join(';'),
      'date': strDate,
      'page': 0,
      'perPage': 30,
  };
  $.post('./commits', params, function(data, textStatus) {
      document.getElementById('buttonList').innerHTML = compiledButton({
          button_data: data
      });
  }, "json");
}
