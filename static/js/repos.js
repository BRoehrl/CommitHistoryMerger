$(".repoSelection").on('change', function(event) {
  var params =  { 'repo': this.name, 'branch': this.value };
  $.post('./repositories', params);
});
