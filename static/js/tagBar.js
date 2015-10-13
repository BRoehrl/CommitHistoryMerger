var isLoading = false
window.onload = function() {
    isLoading = true
    var tags = sessionStorage.getItem('tags');
    if (tags !== null) $(tagBar).val(tags);
    var tagArray = tags.split(",");
    for (var i = 0; i < tagArray.length; i++) {
      $(tagBar).tagsinput('add', tagArray[i])
    }
    isLoading = false
    refreshQuery()
}

$(tagBar).on('beforeItemAdd', function(event) {
  if (event.item.indexOf(':') == -1){
    event.cancel = true
    return
  }
});

$(tagBar).on('itemAdded', function() {
  sessionStorage.setItem('tags', $(tagBar).val());
  if (!isLoading) refreshQuery();
});

$(tagBar).on('itemRemoved', function() {
  sessionStorage.setItem('tags', $(tagBar).val());
  refreshQuery();
});

function refreshQuery(){
  var authors = []
  var repos = []
  var dates = []


  var items = $(tagBar).tagsinput('items')
  var arrayLength = items.length;
  for (var i = 0; i < arrayLength; i++) {
    var wholeTag = items[i];
    var type = wholeTag.substring(0, wholeTag.indexOf(":"));
    var tag = wholeTag.substring(wholeTag.indexOf(":")+1);
    switch (type){
      case 'Author':
      	authors.push(tag)
        break;
      case 'Repo':
        repos.push(tag)
        break;
      case 'Since':
        dates.push(tag)
        break;
      default:
        return;
		}
    }

    var query = './'
    if(authors.length>0){
      query = query + '&author=' + authors.join(';');
    }
    if(repos.length>0){
      query = query +  '&repo=' + repos.join(';');
    }
    if(dates.length>0){
      var earliestDate = new Date(3000, 1, 1);
      var earliestString = '3000-01-01';
      for (var i = 0 ; i < dates.length; i++){
        var dateParts = dates[i].split("-");
        var date = new Date(dateParts[0], (dateParts[1] - 1), dateParts[2]);
        if (date <= earliestDate) {
          earliestDate = date;
          earliestString = dates[i];
          }
      }
      strDate = earliestString + 'T00:00:00Z';
      query = query +  '&since=' + strDate;
    }
    query = query.replace('&','?');
    var oldPath = "." + window.location.href.substring(window.location.href.lastIndexOf("/"));
    if (oldPath != encodeURI(query)) window.location = query;
}
